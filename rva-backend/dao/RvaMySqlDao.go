package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var lockRvaMySqlDao = &sync.Mutex{}
var instanceRvaMySqlDao *RvaMySqlDao

type RvaMySqlDao struct {
	driver     string
	host       string
	port       string
	dbName     string
	user       string
	password   string
	charset    string
	connection *sql.DB
}

func GetRvaMySqlDao() *RvaMySqlDao {
	if instanceRvaMySqlDao == nil {
		lockRvaMySqlDao.Lock()
		defer lockRvaMySqlDao.Unlock()
		if instanceRvaMySqlDao == nil {
			viper.AddConfigPath("./")
			viper.SetConfigName("database")
			viper.ReadInConfig()
			instanceRvaMySqlDao = &RvaMySqlDao{
				driver:   viper.GetString("default.driver"),
				host:     viper.GetString("default.host"),
				port:     viper.GetString("default.port"),
				dbName:   viper.GetString("default.database"),
				user:     viper.GetString("default.user"),
				password: viper.GetString("default.password"),
				charset:  viper.GetString("default.charset"),
			}

		}
	}
	return instanceRvaMySqlDao
}

func (pointer RvaMySqlDao) DeployDatabase() {

	connection, err := sql.Open(instanceRvaMySqlDao.driver, fmt.Sprint(instanceRvaMySqlDao.user, ":", instanceRvaMySqlDao.password, "@tcp(", instanceRvaMySqlDao.host, ":", instanceRvaMySqlDao.port, ")/"))
	if err != nil {
		return
	}
	defer connection.Close()

	connection.Exec("create database if not exists " + pointer.dbName)
	connection.Exec("use " + pointer.dbName)

	connection.Exec(
		`create table if not exists rvaError
	(
		idRvaError int auto_increment,
		message text,
		createdDate datetime default(now()),
		constraint pk_rvaError_idRvaError primary key(idRvaError)
	);`)

	connection.Exec(
		`create table if not exists rvaFileDeployed
	(
		idRvaFileDeployed int auto_increment,
		fileName varchar(200) not null,
		route varchar(700) not null,
		reDeploy boolean default(false),
		creatorAccount varchar(100) not null,
		updaterAccount varchar(100) not null,
		createdDate datetime default(now()),
		updatedDate datetime default(now()),
		constraint pk_rvaFileDeployed_idRvaFileDeployed primary key(idRvaFileDeployed),
		constraint uk_rvaFileDeployed_route unique key(route)
	);`)

	deployByFolder := func(folder string) error {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		basepath = filepath.Clean(filepath.Join(basepath, ".."))
		basepath = filepath.Clean(filepath.Join(basepath, ".."))
		basepath = basepath + "\\rva-database\\" + folder

		files, err := ioutil.ReadDir(basepath)
		if err != nil {
			return err
		}

		context := context.Background()
		transaction, err := connection.BeginTx(context, nil)
		if err != nil {
			transaction.Rollback()
			return err
		}
		transaction.ExecContext(context, "set autocommit=0;")

		for _, file := range files {

			route := basepath + "\\" + file.Name()
			fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			isAlreadyDeployed := false
			reDeploy := false

			transaction.QueryRowContext(context, "select if(reDeploy,false,if(count(*)>0,true,false)) from rvaFileDeployed where route = ? limit 1", route).Scan(&isAlreadyDeployed)

			if !isAlreadyDeployed {
				fileBytes, _ := ioutil.ReadFile(route)

				if folder == "procedure" || folder == "trigger" {
					reDeploy = true
					transaction.ExecContext(context, "drop "+folder+" if exists "+fileName+";")
				}

				_, err = transaction.ExecContext(context, string(fileBytes))
				if err != nil {
					transaction.Rollback()
					return errors.New("DEPLOYMENT-ERROR in " + route + " -> " + err.Error())
				}

				if folder == "procedure" {
					var procedureQuery string
					err = transaction.QueryRowContext(context, "select concat('call ',?,'(',ifnull(group_concat(PARAMETER_NAME separator ','),''),');') from information_schema.parameters where SPECIFIC_SCHEMA = ? and SPECIFIC_NAME = ? and ROUTINE_TYPE='PROCEDURE' limit 1", fileName, pointer.dbName, fileName).Scan(&procedureQuery)
					if err != nil {
						transaction.Rollback()
						return err
					}
					_, err = transaction.ExecContext(context, "insert into rvaProcedure (procedureName,procedureQuery,creatorAccount,updaterAccount) values (?,?,?,?);", fileName, procedureQuery, "System", "System")
					if err != nil {
						transaction.ExecContext(context, "update rvaProcedure set procedureQuery = ?, updaterAccount = ?, updatedDate = now() where procedureName = ?;", procedureQuery, "System", fileName)
					}
				}

				_, err = transaction.ExecContext(context, "insert into rvaFileDeployed (fileName,route,reDeploy,creatorAccount,updaterAccount) values (?,?,?,?,?);", fileName, route, reDeploy, "System", "System")
				if err != nil {
					transaction.ExecContext(context, "update rvaFileDeployed set updaterAccount = ?, updatedDate = now() where route = ?;", "System", route)
				}
			}
		}
		transaction.Commit()

		return nil
	}

	for _, value := range []string{"schema", "function", "view", "procedure", "event", "trigger", "data", "migration"} {
		err = deployByFolder(value)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}

func (pointer RvaMySqlDao) OpenConnection() *sql.DB {
	pointer.connection, _ = sql.Open(instanceRvaMySqlDao.driver, fmt.Sprint(instanceRvaMySqlDao.user, ":", instanceRvaMySqlDao.password, "@tcp(", instanceRvaMySqlDao.host, ":", instanceRvaMySqlDao.port, ")/", instanceRvaMySqlDao.dbName, "?charset=", instanceRvaMySqlDao.charset, "&parseTime=True&loc=Local"))
	return pointer.connection
}

func (pointer RvaMySqlDao) Execute(query string) interface{} {
	return pointer.execute(query, nil)
}

func (pointer RvaMySqlDao) ExecuteWithoutLock(query string) interface{} {
	return pointer.execute(query, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
}

func (pointer RvaMySqlDao) ExecuteContext(parameter interface{}, queries []string) (interface{}, error) {
	connection := pointer.OpenConnection()
	context := context.Background()
	transaction, err := connection.BeginTx(context, nil)
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	executeQuery := func(query string, parameter map[string]interface{}) (interface{}, error) {
		rows, err := transaction.QueryContext(context, pointer.mapParameters(query, parameter))
		defer rows.Close()
		if err != nil {
			transaction.Rollback()
			return nil, err
		}
		return pointer.mapResultset(rows), nil
	}

	var finalResult interface{}

	for index, query := range queries {
		var recursive func(parameter interface{}) (interface{}, error)

		recursive = func(parameter interface{}) (interface{}, error) {
			var result interface{}
			var results []interface{}
			switch parameter.(type) {
			case map[string]interface{}:
				return executeQuery(query, parameter.(map[string]interface{}))
			case []map[string]interface{}:
				for _, value := range parameter.([]map[string]interface{}) {
					result, err = executeQuery(query, value)
					if err != nil {
						return nil, err
					}
					results = append(results, result)
				}
				return results, nil
			case []interface{}:
				for _, parameter := range parameter.([]interface{}) {
					result, err = recursive(parameter)
					if err != nil {
						return nil, err
					}
					results = append(results, result)
				}
			}
			return results, nil
		}

		result,err:=recursive(parameter)
		if err != nil {
			return nil, err
		}
		if index == len(queries)-1 {
			finalResult = result
		} else {
			parameter = result
		}
	}
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	return finalResult, nil
}

func (pointer RvaMySqlDao) execute(query string, opts *sql.TxOptions) interface{} {
	connection := pointer.OpenConnection()
	context := context.Background()
	transaction, err := connection.BeginTx(context, opts)
	if err != nil {
		transaction.Rollback()
		return nil
	}
	rows, _ := transaction.QueryContext(context, query)
	defer rows.Close()
	if err != nil {
		transaction.Rollback()
		return nil
	}
	result := pointer.mapResultset(rows)
	err = transaction.Commit()
	if err != nil {
		transaction.Rollback()
		return nil
	}

	return result
}

func (pointer RvaMySqlDao) mapParameters(query string, values map[string]interface{}) string {
	parameters := strings.Split(strings.Split(strings.Split(query, "(")[1], ")")[0], ",")
	for index, parameter := range parameters {
		if parameter == "" {
			parameters = append(parameters[:index], parameters[index+1:]...)
			break
		}
	}
	parameterFirm := ""
	if len(parameters) > 0 {
		for _, parameter := range parameters {
			if values[parameter] != nil {
				switch values[parameter].(type) {
				default:
					parameterFirm = parameterFirm + fmt.Sprintf("%v,", values[parameter])
				case string:
					parameterFirm = parameterFirm + fmt.Sprintf("'%v',", values[parameter])
				case time.Time:
					parameterFirm = parameterFirm + fmt.Sprintf("'%v',", values[parameter].(time.Time).Format(time.RFC3339))
				}
			} else {
				parameterFirm = parameterFirm + "null,"
			}
		}
		return strings.Split(query, "(")[0] + "(" + parameterFirm[:len(parameterFirm)-1] + ");"
	}
	return query
}

func (pointer RvaMySqlDao) mapResultset(rows *sql.Rows) interface{} {
	mapper := func() interface{} {
		columns, _ := rows.ColumnTypes()
		count := len(columns)
		dataset := make([]map[string]interface{}, 0)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)

		countRows := 0
		for rows.Next() {

			for index := 0; index < count; index++ {
				valuePtrs[index] = &values[index]
			}

			rows.Scan(valuePtrs...)
			datarow := make(map[string]interface{})
			for index, col := range columns {
				var v interface{}
				val := values[index]
				b, ok := val.([]byte)
				if ok {
					v = string(b)
					switch strings.ToUpper(col.DatabaseTypeName()) {
					case "INT":
						v, _ = strconv.Atoi(v.(string))
					}

				} else {
					v = val
				}
				datarow[col.Name()] = v
			}

			dataset = append(dataset, datarow)
			countRows++
		}

		return func() interface{} {
			if countRows == 1 {
				return dataset[0]
			} else if countRows > 1 {
				return dataset
			} else {
				return nil
			}
		}()
	}

	var dataset []interface{}
	countTables := 0
	for {
		dataset = append(dataset, mapper())
		countTables++
		if !rows.NextResultSet() {
			break
		}
	}

	return func() interface{} {
		if countTables > 1 {
			return dataset
		} else {
			return dataset[0]
		}
	}()
}
