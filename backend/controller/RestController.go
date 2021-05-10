package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rva/client"
	"rva/helper"
	"rva/strategy"
	"strings"
	"sync"
)

var multitonRestController = make(map[interface{}]*RestController)
var onceRestController = make(map[interface{}]*sync.Once)

type RestController struct {
	jwtClient    *client.JwtClient
	recursionHelper *helper.RecursionHelper
	mapperHelper *helper.MapperHelper
}

func GetRestController(jwtClient *client.JwtClient) *RestController {
	return &RestController{
		jwtClient:    jwtClient,
		recursionHelper: helper.GetRecursionHelper(),
		mapperHelper: helper.GetMapperHelper(),
	}
}

func (pointer RestController) GetRouter() http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())

	buildEndpoint := func(endpoint map[string]interface{}) {
		router.Handle(strings.ToUpper(endpoint["httpVerb"].(string)), endpoint["path"].(string), func(context *gin.Context) {
			var request interface{}
			var statusCode int
			var response interface{}
			var isValid bool

			context.ShouldBindJSON(&request)

			if request == nil {
				request = make(map[string]interface{})
			}
			for _, value := range context.Params {
				pointer.recursionHelper.RecursionOfFunctionWithMapParameterAndWithoutReturn(func(parameter map[string]interface{}) {parameter[value.Key] = value.Value},request)
			}
			for key, value := range context.Request.URL.Query() {
				if len(value) > 0 { 
					pointer.recursionHelper.RecursionOfFunctionWithMapParameterAndWithoutReturn(func(parameter map[string]interface{}) {parameter[key] = value[0]},request)
				}
			}
			if endpoint["validAuth"].(bool) {
				if context.GetHeader("Authorization") == "" {
					isValid, statusCode, response = false, http.StatusBadRequest, map[string]interface{}{
						"Status":  false,
						"Message": "No token found",
					}
				}
				if pointer.jwtClient.ValidateToken(context.GetHeader("Authorization")) {
					isValid = true
				}
				isValid, statusCode, response = false, http.StatusUnauthorized, map[string]interface{}{
					"Status":  false,
					"Message": "Token is not valid",
				}
			} else {
				isValid = true
			}

			if isValid || !endpoint["validAuth"].(bool) {
				response = pointer.sqlDao.ExecuteProcedureGroup(endpoint["idProcedureGroup"].(int), request)
			}

			if endpoint["createAuth"].(bool) && response != nil {
				context.Header("Authorization", pointer.jwtClient.GenerateToken(response.(map[string]interface{})))
			}

			if statusCode == 0 {
				if response != nil {
					statusCode = http.StatusOK
				} else {
					statusCode = http.StatusNoContent
				}
			}

			if statusCode == 0 {
				if response != nil {
					statusCode = endpoint["statusCodeSuccess"].(int)
				} else {
					statusCode = endpoint["statusCodeNoResponse"].(int)
				}
			}
			if endpoint["retriveOnlyHeaders"].(bool) {
				function := func(parameter map[string]interface{}) {
					for key, value := range parameter {
						context.Header(key, fmt.Sprint(value))
					}
				}
				pointer.mapperHelper.ParameterizedRecursionNotReturned(function, response)
				context.Status(statusCode)
			} else {
				context.JSON(statusCode, response)
			}
		})
	}

	pointer.mapperHelper.ParameterizedRecursionNotReturned(buildEndpoint, pointer.sqlDao.GetEndpoints())

	return router
}
