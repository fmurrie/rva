package sql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"rva/helper"
	"strconv"
	"strings"
	"sync"
	"time"
)

var instanceMySqlDao = make(map[interface{}]*MySqlDao)
var singletonMySqlDao = make(map[interface{}]*sync.Once)

type MySqlDao struct {
	database     ISqlDatabase
	connection   *sql.DB
	mapperHelper *helper.MapperHelper
}

func GetMySqlDao(database ISqlDatabase) *MySqlDao {
	if singletonMySqlDao[database] == nil {
		var once sync.Once
		singletonMySqlDao[database] = &once
	}
	singletonMySqlDao[database].Do(func() {
		connection, _ := sql.Open(database.GetDriver(), fmt.Sprint(database.GetUser(), ":", database.GetPassword(), "@tcp(", database.GetHost(), ":", database.GetPort(), ")/", database.GetDbName(), "?charset=", database.GetCharset(), "&parseTime=True&loc=Local"))
		instanceMySqlDao[database] = &MySqlDao{
			database:     database,
			connection:   connection,
			mapperHelper: helper.GetMapperHelper(),
		}
	})
	return instanceMySqlDao[database]
}

func (pointer MySqlDao) MapSqlQueryParameters(query string, values map[string]interface{}) string {
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

func (pointer MySqlDao) MapSqlRows(rows *sql.Rows) interface{} {
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
					case "TINYINT":
						v, _ = strconv.ParseBool(v.(string))
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

func (pointer MySqlDao) ExecuteProcedureGroup(idProcedureGroup int, parameter interface{}) interface{} {
	var finalResult interface{}
	var queries []string
	function := func(parameter map[string]interface{}) {
		queries = append(queries, parameter["query"].(string))
	}
	pointer.mapperHelper.ParameterizedRecursionNotReturned(function, pointer.GetProcedureGroupSteps(idProcedureGroup))

	context := context.Background()
	transaction, err := pointer.connection.BeginTx(context, nil)
	if err != nil {
		transaction.Rollback()
		pointer.LogError(err)
		return nil
	}

	for index, query := range queries {
		function := func(parameter map[string]interface{}) (interface{}, error) {
			rows, err := transaction.QueryContext(context, pointer.MapSqlQueryParameters(query, parameter))
			defer rows.Close()
			if err != nil {
				transaction.Rollback()
				pointer.LogError(err)
				return nil, err
			}
			return pointer.MapSqlRows(rows), nil
		}
		result, err := pointer.mapperHelper.ParameterizedRecursionReturnedWithError(function, parameter)
		if err != nil {
			return nil
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
		pointer.LogError(err)
		return nil
	}

	return finalResult
}

func (pointer MySqlDao) GetEndpoints() interface{} {
	context := context.Background()
	transaction, _ := pointer.connection.BeginTx(context, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	rows, _ := transaction.QueryContext(context, "call sys_endpoint_getAll();")
	defer rows.Close()
	result := pointer.MapSqlRows(rows)
	transaction.Commit()
	return result
}

func (pointer MySqlDao) GetProcedureGroupSteps(idProcedureGroup int) interface{} {
	context := context.Background()
	transaction, _ := pointer.connection.BeginTx(context, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	rows, _ := transaction.QueryContext(context, fmt.Sprint("call sys_procedureGroupStep_getByIdProcedureGroup(", idProcedureGroup, ");"))
	defer rows.Close()
	result := pointer.MapSqlRows(rows)
	transaction.Commit()
	return result
}

func (pointer MySqlDao) LogError(err error) {
	pointer.connection.Exec("insert into sys_error (message) values (?);", err.Error())
}
