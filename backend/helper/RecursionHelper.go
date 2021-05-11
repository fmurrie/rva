package helper

import (
	"sync"
)

var singletonRecursionHelper *RecursionHelper
var onceRecursionHelper sync.Once

type RecursionHelper struct {
}

func GetRecursionHelper() *RecursionHelper {
	onceRecursionHelper.Do(func() {
		singletonRecursionHelper = &RecursionHelper{}
	})
	return singletonRecursionHelper
}

func (pointer RecursionHelper) RecursionOfFunctionWithMapParameterAndInterfaceReturn(function func(parameter map[string]interface{}) interface{}, parameters ...interface{}) interface{} {

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
				results = append(results, pointer.RecursionOfFunctionWithMapParameterAndInterfaceReturn(function, value))
			}
		}
		if len(parameters) <= 1 {
			return results
		}
		resultsGroup = append(resultsGroup, results)
	}

	return resultsGroup
}

func (pointer RecursionHelper) RecursionOfFunctionWithMapParameterAndInterfaceReturnOrGiveError(function func(parameter map[string]interface{}) (interface{},error), parameters ...interface{}) (interface{},error) {

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
				result,err :=pointer.RecursionOfFunctionWithMapParameterAndInterfaceReturnOrGiveError(function, value)
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

func (pointer RecursionHelper) RecursionOfFunctionWithMapParameterAndWithoutReturn(function func(parameter map[string]interface{}), parameters ...interface{}) {
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
				pointer.RecursionOfFunctionWithMapParameterAndWithoutReturn(function, value)
			}
		}
	}
}

func (pointer RecursionHelper) SetKeyAndValueInAnyMapInsideTheInterfaceParameter(parameter interface{}, paramKey string, paramValue interface{}) {
	function := func(parameter map[string]interface{}) {
		parameter[paramKey] = paramValue
	}
	pointer.RecursionOfFunctionWithMapParameterAndWithoutReturn(function, parameter)
}