package sql

import "database/sql"

type ISqlDao interface {
	MapSqlQueryParameters(query string, values map[string]interface{}) string
	MapSqlRows(rows *sql.Rows) interface{}
	ExecuteProcedureGroup(idProcedureGroup int, parameter interface{}) interface{}
	GetEndpoints() interface{}
	GetProcedureGroupSteps(idProcedureGroup int) interface{}
	LogError(err error)
}
