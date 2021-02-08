package auth

import (
	"rva/keyword"
	"sync"
)

var instanceAuthControllerFactory *AuthControllerFactory
var singletonAuthControllerFactory sync.Once

type AuthControllerFactory struct {
}

func GetAuthControllerFactory() *AuthControllerFactory {
	singletonAuthControllerFactory.Do(func() {
		instanceAuthControllerFactory = &AuthControllerFactory{}
	})
	return instanceAuthControllerFactory
}

func (pointer AuthControllerFactory) GetAuthController(key string) IAuthController {
	switch key {
	default:
		return nil
	case keyword.Jwt:
		return GetJwtAuthController()
	}
}
