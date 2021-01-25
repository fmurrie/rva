package main

import (
	"rva/dao"
	"rva/controller"
)



func main() {
	dao.GetRvaMySqlDao().DeployDatabase()
	controller.GetRvaRestController().RaiseEndpoints()
}
