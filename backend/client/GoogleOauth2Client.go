package client

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"rva/helper"
	"sync"
)

var multitonGoogleOauth2Client = make(map[interface{}]*GoogleOauth2Client)
var onceGoogleOauth2Client = make(map[interface{}]*sync.Once)

type GoogleOauth2Client struct {
	config *oauth2.Config
	responseURL string
}

func GetGoogleOauth2Client(configurationHelper *helper.ConfigurationHelper) *GoogleOauth2Client {
	var locker sync.Mutex
	locker.Lock()
    defer locker.Unlock()
	if onceGoogleOauth2Client[configurationHelper] == nil {
		onceGoogleOauth2Client[configurationHelper] =  &sync.Once{}
	}
	onceGoogleOauth2Client[configurationHelper].Do(func() {
		securityHelper := helper.GetSecurityHelper()
		multitonGoogleOauth2Client[configurationHelper] = &GoogleOauth2Client{
			config: &oauth2.Config{
				ClientID:     securityHelper.Decrypt(configurationHelper.GetStringValueByKey("client_id")),
				ClientSecret: securityHelper.Decrypt(configurationHelper.GetStringValueByKey("client_secret")),
				RedirectURL:  securityHelper.Decrypt(configurationHelper.GetStringValueByKey("callback_path")),
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
				},
				Endpoint: google.Endpoint,
			},
			responseURL:securityHelper.Decrypt(configurationHelper.GetStringValueByKey("response_url")),
		}
	})
	return multitonGoogleOauth2Client[configurationHelper]
}

func (pointer GoogleOauth2Client) GetConfig() *oauth2.Config{
	return pointer.config
}

func (pointer GoogleOauth2Client) GetResponseURL() string{
	return pointer.responseURL
}