package controller

import (
	"fmt"
	"rva/factory"
	"rva/helper"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var lockRvaAuthController = &sync.Mutex{}
var instanceRvaAuthController *RvaAuthController

type RvaAuthController struct {
	signinKey string
	rvaDao factory.RvaDao
}

func GetRvaAuthController() *RvaAuthController {
	if instanceRvaAuthController == nil {
		lockRvaAuthController.Lock()
		defer lockRvaAuthController.Unlock()
		if instanceRvaAuthController == nil {
			viper.AddConfigPath("./")
			viper.SetConfigName("configuration")
			viper.ReadInConfig()
			instanceRvaAuthController = &RvaAuthController{
				signinKey: helper.GetRvaSecurityHelper().Decrypt(viper.GetString("authentication.signin_key")),
				rvaDao: factory.GetRvaDao(),
			}
		}
	}
	return instanceRvaAuthController
}

func (pointer RvaAuthController) GenerateToken(account map[string]interface{}){
	if account != nil {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["authorized"] = true
		claims["account"] = account
		claims["issuedAt"] = time.Now().Unix()
		claims["expiresAt"] = time.Now().Add(time.Minute * 15).Unix()
		tokenString, err := token.SignedString([]byte(pointer.signinKey))
		if err != nil {
			pointer.rvaDao.LogError(err)
		} else {
			account["token"] = tokenString
		}
	}
}

func (pointer RvaAuthController) ValidateToken(token string) bool {
	result, err := jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(pointer.signinKey), nil
	})

	if err != nil {
		pointer.rvaDao.LogError(err)
		return false
	}

	return result.Valid
}

func (pointer RvaAuthController) GetClaims(token *jwt.Token) jwt.MapClaims {
	return token.Claims.(jwt.MapClaims)
}
