package main

import (
	"rva/client"
	"rva/controller"
	"rva/controller/auth"
	"rva/dao/sql"
	"rva/keyword"
)

func main() {
	client.GetDatabaseClient().DeployDatabase()
	controller.GetRestController(sql.GetMySqlDao(client.GetDatabaseClient()), auth.GetAuthControllerFactory().GetAuthController(keyword.Jwt)).RaiseEndpoints()
}