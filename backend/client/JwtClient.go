package client

import (
	"fmt"
	"rva/helper"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var multitonJwtClient = make(map[interface{}]*JwtClient)
var onceJwtClient = make(map[interface{}]*sync.Once)

type JwtClient struct {
	signinKey string
}

func GetJwtClient(configurationHelper *helper.ConfigurationHelper) *JwtClient {
	var locker sync.Mutex
	locker.Lock()
    defer locker.Unlock()
	if onceJwtClient[configurationHelper] == nil {
		onceJwtClient[configurationHelper] = &sync.Once{}
	}
	onceJwtClient[configurationHelper].Do(func() {
		securityHelper := helper.GetSecurityHelper()
		multitonJwtClient[configurationHelper] = &JwtClient{
			signinKey: securityHelper.Decrypt(configurationHelper.GetStringValueByKey("signin_key")),
		}
	})
	return multitonJwtClient[configurationHelper]
}

func (pointer JwtClient) GenerateToken(information interface{}) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["information"] = information
	claims["issuedAt"] = time.Now().Unix()
	claims["expiresAt"] = time.Now().Add(time.Minute * 15).Unix()
	authToken,_:=token.SignedString([]byte(pointer.signinKey))
	return authToken
}

func (pointer JwtClient) ValidateToken(token string) bool {
result, err := jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
	}
	return []byte(pointer.signinKey), nil
})
if err != nil {
	return false
}
return result.Valid
}

func (pointer JwtClient) GetClaims(token *jwt.Token) jwt.MapClaims {
return token.Claims.(jwt.MapClaims)
}