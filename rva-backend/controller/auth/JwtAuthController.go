package auth

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"rva/client"
	"sync"
	"time"
)

var instanceJwtAuthController *JwtAuthController
var singletonJwtAuthController sync.Once

type JwtAuthController struct {
	jwtClient *client.JwtClient
}

func GetJwtAuthController() *JwtAuthController {
	singletonJwtAuthController.Do(func() {
		instanceJwtAuthController = &JwtAuthController{
			jwtClient: client.GetJwtClient(),
		}
	})
	return instanceJwtAuthController
}

func (pointer JwtAuthController) GenerateToken(information interface{}) string {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["authorized"] = true
		claims["information"] = information
		claims["issuedAt"] = time.Now().Unix()
		claims["expiresAt"] = time.Now().Add(time.Minute * 15).Unix()
		authToken,_:=token.SignedString([]byte(pointer.jwtClient.GetSigninKey()))
		return authToken
}

func (pointer JwtAuthController) ValidateToken(token string) bool {
	result, err := jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(pointer.jwtClient.GetSigninKey()), nil
	})
	if err != nil {
		return false
	}
	return result.Valid
}

func (pointer JwtAuthController) GetClaims(token *jwt.Token) jwt.MapClaims {
	return token.Claims.(jwt.MapClaims)
}
