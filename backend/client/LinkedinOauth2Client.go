package client

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
	"rva/helper"
	"sync"
)

var multitonLinkedinOauth2Client = make(map[interface{}]*LinkedinOauth2Client)
var onceLinkedinOauth2Client = make(map[interface{}]*sync.Once)

type LinkedinOauth2Client struct {
	config *oauth2.Config
	responseURL string
}

func GetLinkedinOauth2Client(configurationHelper *helper.ConfigurationHelper) *LinkedinOauth2Client {
	if onceLinkedinOauth2Client[configurationHelper] == nil {
		onceLinkedinOauth2Client[configurationHelper] = &sync.Once{}
	}
	onceLinkedinOauth2Client[configurationHelper].Do(func() {
		securityHelper := helper.GetSecurityHelper()
		multitonLinkedinOauth2Client[configurationHelper] = &LinkedinOauth2Client{
			config: &oauth2.Config{
				ClientID:     securityHelper.Decrypt(configurationHelper.GetStringValueByKey("client_id")),
				ClientSecret: securityHelper.Decrypt(configurationHelper.GetStringValueByKey("client_secret")),
				RedirectURL:  securityHelper.Decrypt(configurationHelper.GetStringValueByKey("callback_path")),
				Scopes: []string{
					"r_basicprofile",
					"r_emailaddress",
				},
				Endpoint: linkedin.Endpoint,
			},
			responseURL:securityHelper.Decrypt(configurationHelper.GetStringValueByKey("response_url")),
		}
	})
	return multitonLinkedinOauth2Client[configurationHelper]
}

func (pointer LinkedinOauth2Client) GetConfig() *oauth2.Config{
	return pointer.config
}

func (pointer LinkedinOauth2Client) GetResponseURL() string{
	return pointer.responseURL
}