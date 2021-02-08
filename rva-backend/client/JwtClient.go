package client

import (
	"rva/helper"
	"sync"
)

var instanceJwtClient *JwtClient
var singletonJwtClient sync.Once

type JwtClient struct {
	signinKey string
}

func GetJwtClient() *JwtClient {
	singletonJwtClient.Do(func() {
		configurationHelper := helper.GetConfigurationHelper("jwt", "configuration", "./")
		securityHelper := helper.GetSecurityHelper()
		instanceJwtClient = &JwtClient{
			signinKey: securityHelper.Decrypt(configurationHelper.GetStringValueByKey("signin_key")),
		}
	})
	return instanceJwtClient
}

func (pointer JwtClient) GetSigninKey() string {
	return pointer.signinKey
}
