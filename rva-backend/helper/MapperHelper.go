package helper

import (
	"sync"
)

var instanceMapperHelper *MapperHelper
var singletonMapperHelper sync.Once

type MapperHelper struct {
}

func GetMapperHelper() *MapperHelper {
	singletonMapperHelper.Do(func() {
		instanceMapperHelper = &MapperHelper{}
	})
	return instanceMapperHelper
}

func (pointer MapperHelper) ParameterizedRecursionReturned(function func(parameter map[string]interface{}) interface{}, parameters ...interface{}) interface{} {

	var resultsGroup []interface{}
	for _, parameter := range parameters {
		var results []interface{}

		switch parameter.(type) {
		case map[string]interface{}:
			result := function(parameter.(map[string]interface{}))
			if len(parameters) <= 1 {
				return result
			}
			results = append(results, result)
		case []map[string]interface{}:
			for _, value := range parameter.([]map[string]interface{}) {
				results = append(results, function(value))
			}
			return results
		case []interface{}:
			for _, value := range parameter.([]interface{}) {
				results = append(results, pointer.ParameterizedRecursionReturned(function, value))
			}
		}
		if len(parameters) <= 1 {
			return results
		}
		resultsGroup = append(resultsGroup, results)
	}

	return resultsGroup
}

func (pointer MapperHelper) ParameterizedRecursionReturnedWithError(function func(parameter map[string]interface{}) (interface{},error), parameters ...interface{}) (interface{},error) {

	var resultsGroup []interface{}
	for _, parameter := range parameters {
		var results []interface{}

		switch parameter.(type) {
		case map[string]interface{}:
			result,err := function(parameter.(map[string]interface{}))
			if len(parameters) <= 1 || err!=nil{
				return result,err
			}
			results = append(results, result)
		case []map[string]interface{}:
			for _, value := range parameter.([]map[string]interface{}) {
				result,err := function(value)
				if err!=nil{
					return results,err
				}
				results = append(results, result)
			}
			return results,nil
		case []interface{}:
			for _, value := range parameter.([]interface{}) {
				result,err :=pointer.ParameterizedRecursionReturnedWithError(function, value)
				if err!=nil{
					return results,err
				}
				results = append(results, result)
			}
		}
		if len(parameters) <= 1 {
			return results,nil
		}
		resultsGroup = append(resultsGroup, results)
	}

	return resultsGroup,nil
}

func (pointer MapperHelper) ParameterizedRecursionNotReturned(function func(parameter map[string]interface{}), parameters ...interface{}) {
	for _, parameter := range parameters {
		switch parameter.(type) {
		case map[string]interface{}:
			function(parameter.(map[string]interface{}))
		case []map[string]interface{}:
			for _, value := range parameter.([]map[string]interface{}) {
				function(value)
			}
		case []interface{}:
			for _, value := range parameter.([]interface{}) {
				pointer.ParameterizedRecursionNotReturned(function, value)
			}
		}
	}
}

func (pointer MapperHelper) MapKeyValue(parameter interface{}, paramKey string, paramValue interface{}) {
	function := func(parameter map[string]interface{}) {
		parameter[paramKey] = paramValue
	}
	pointer.ParameterizedRecursionNotReturned(function, parameter)
}