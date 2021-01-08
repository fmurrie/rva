package main

import (
	"rva/controller"
	"rva/dao"
)

func main() {
	dao.GetRvaMySqlDao().DeployDatabase()
	controller.GetRvaRestController().RaiseEndpoints()
}
