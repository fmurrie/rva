package client

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"rva/helper"
	"strings"
	"sync"
)

var instanceDatabaseClient *DatabaseClient
var singletonDatabaseClient sync.Once

type DatabaseClient struct {
	driver   string
	host     string
	port     string
	dbName   string
	user     string
	password string
	charset  string
	deploy   bool
}

func GetDatabaseClient() *DatabaseClient {
	singletonDatabaseClient.Do(func() {
		configurationHelper := helper.GetConfigurationHelper("database", "configuration", "./")
		securityHelper := helper.GetSecurityHelper()
		instanceDatabaseClient = &DatabaseClient{
			driver:   securityHelper.Decrypt(configurationHelper.GetStringValueByKey("driver")),
			host:     securityHelper.Decrypt(configurationHelper.GetStringValueByKey("host")),
			port:     securityHelper.Decrypt(configurationHelper.GetStringValueByKey("port")),
			dbName:   securityHelper.Decrypt(configurationHelper.GetStringValueByKey("dbname")),
			user:     securityHelper.Decrypt(configurationHelper.GetStringValueByKey("user")),
			password: securityHelper.Decrypt(configurationHelper.GetStringValueByKey("password")),
			charset:  securityHelper.Decrypt(configurationHelper.GetStringValueByKey("charset")),
			deploy:   configurationHelper.GetBoolValueByKey("deploy"),
		}
	})
	return instanceDatabaseClient
}

func (pointer DatabaseClient) GetDriver() string {
	return pointer.driver
}

func (pointer DatabaseClient) GetHost() string {
	return pointer.host
}

func (pointer DatabaseClient) GetPort() string {
	return pointer.port
}

func (pointer DatabaseClient) GetDbName() string {
	return pointer.dbName
}

func (pointer DatabaseClient) GetUser() string {
	return pointer.user
}

func (pointer DatabaseClient) GetPassword() string {
	return pointer.password
}

func (pointer DatabaseClient) GetCharset() string {
	return pointer.charset
}

func (pointer DatabaseClient) DeployDatabase() {
	if pointer.deploy {
		connection, err := sql.Open(pointer.GetDriver(), fmt.Sprint(pointer.GetUser(), ":", pointer.GetPassword(), "@tcp(", pointer.GetHost(), ":", pointer.GetPort(), ")/"))
		if err != nil {
			return
		}
		defer connection.Close()

		connection.Exec("create database if not exists " + pointer.GetDbName())
		connection.Exec("use " + pointer.GetDbName())
		connection.Exec(
			`create table if not exists sys_error
		(
			idError int auto_increment,
			message text,
			createdDate datetime default(now()),
			constraint pk_sys_error_idError primary key(idError)
		);`)

		connection.Exec(
			`create table if not exists sys_fileDeployed
		(
			idFileDeployed int auto_increment,
			fileName varchar(200) not null,
			route varchar(700) not null,
			reDeploy boolean default(false),
			creatorIp varchar(100) not null,
			updaterIp varchar(100) not null,
			createdDate datetime default(now()),
			updatedDate datetime default(now()),
			constraint pk_sys_fileDeployed_idFileDeployed primary key(idFileDeployed),
			constraint uk_sys_fileDeployed_route unique key(route)
		);`)

		deployByFolder := func(folder string) error {
			_, b, _, _ := runtime.Caller(0)
			basepath := filepath.Dir(b)
			basepath = filepath.Clean(filepath.Join(basepath, ".."))
			basepath = filepath.Clean(filepath.Join(basepath, ".."))
			basepath = basepath + "\\rva-database\\mysql\\" + folder

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

			for _, file := range files {
				route := basepath + "\\" + file.Name()
				fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
				isAlreadyDeployed := false
				reDeploy := false

				transaction.QueryRowContext(context, "select if(count(*)>0,true,false),ifnull(reDeploy,false) from sys_fileDeployed where route = ? limit 1", route).Scan(&isAlreadyDeployed, &reDeploy)

				if !isAlreadyDeployed || reDeploy {
					fileBytes, _ := ioutil.ReadFile(route)

					if folder == "function" || folder == "view" || folder == "event" || folder == "procedure" || folder == "trigger" {
						reDeploy = true
						transaction.ExecContext(context, "drop "+folder+" if exists "+fileName+";")
						_, err = transaction.ExecContext(context, string(fileBytes))
						if err != nil {
							transaction.Rollback()
							return errors.New("DEPLOYMENT-ERROR in " + route + " -> " + err.Error())
						}
						if folder == "procedure" {
							if !isAlreadyDeployed {
								transaction.ExecContext(context, "insert into sys_procedure (name,query,creatorIp,updaterIp) select SPECIFIC_NAME as name,concat('call ',SPECIFIC_NAME,'(',ifnull(group_concat(PARAMETER_NAME separator ','),''),');') as query, 'System' as creatorIp, 'System' as updaterIp from information_schema.parameters where SPECIFIC_SCHEMA = ? and SPECIFIC_NAME = ? and ROUTINE_TYPE='PROCEDURE' limit 1;", pointer.GetDbName(), fileName)
							} else {
								transaction.ExecContext(context, "update sys_procedure set query = (select concat('call ',SPECIFIC_NAME,'(',ifnull(group_concat(PARAMETER_NAME separator ','),''),');') as query from information_schema.parameters where SPECIFIC_SCHEMA = ? and SPECIFIC_NAME = name and ROUTINE_TYPE='PROCEDURE' limit 1), updaterIp = 'System', updatedDate = now() where name = ? and query != (select concat('call ',SPECIFIC_NAME,'(',ifnull(group_concat(PARAMETER_NAME separator ','),''),');') as query from information_schema.parameters where SPECIFIC_SCHEMA = ? and SPECIFIC_NAME = name and ROUTINE_TYPE='PROCEDURE' limit 1);", pointer.GetDbName(), fileName, pointer.GetDbName())
							}
						}
					} else {
						instructions := strings.Split(string(fileBytes), ";")
						for _, instruction := range instructions {
							if instruction != "" {
								_, err = transaction.ExecContext(context, instruction)
								if err != nil {
									transaction.Rollback()
									return errors.New("DEPLOYMENT-ERROR in " + route + " -> " + err.Error())
								}
							}
						}
					}
					if !isAlreadyDeployed {
						_, err = transaction.ExecContext(context, "insert into sys_fileDeployed (fileName,route,reDeploy,creatorIp,updaterIp) values (?,?,?,?,?);", fileName, route, reDeploy, "System", "System")
					} else {
						transaction.ExecContext(context, "update sys_fileDeployed set updaterIp = ?, updatedDate = now() where route = ?;", "System", route)
					}
				}
			}
			transaction.Commit()
			return nil
		}

		for _, value := range []string{"schema", "function", "view", "procedure", "event", "trigger", "data", "migration"} {
			err := deployByFolder(value)
			if err != nil {
				log.Fatalln(err.Error())
			}
		}
	}
}
