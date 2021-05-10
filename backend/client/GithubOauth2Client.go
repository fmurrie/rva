package client

import (
	"rva/helper"
	"sync"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var multitonGithubOauth2Client = make(map[interface{}]*GithubOauth2Client)
var onceGithubOauth2Client = make(map[interface{}]*sync.Once)

type GithubOauth2Client struct {
	config *oauth2.Config
	responseURL string
}

func GetGithubOauth2Client(configurationHelper *helper.ConfigurationHelper) *GithubOauth2Client {
	var locker sync.Mutex
	locker.Lock()
    defer locker.Unlock()
	if onceGithubOauth2Client[configurationHelper] == nil {
		onceGithubOauth2Client[configurationHelper] = &sync.Once{}
	}
	onceGithubOauth2Client[configurationHelper].Do(func() {
		securityHelper := helper.GetSecurityHelper()
		multitonGithubOauth2Client[configurationHelper] = &GithubOauth2Client{
			config: &oauth2.Config{
				ClientID:     securityHelper.Decrypt(configurationHelper.GetStringValueByKey("client_id")),
				ClientSecret: securityHelper.Decrypt(configurationHelper.GetStringValueByKey("client_secret")),
				RedirectURL:  securityHelper.Decrypt(configurationHelper.GetStringValueByKey("callback_path")),
				Scopes: []string{
					"user:email",
				},
				Endpoint: github.Endpoint,
			},
			responseURL:securityHelper.Decrypt(configurationHelper.GetStringValueByKey("response_url")),
		}
	})
	return multitonGithubOauth2Client[configurationHelper]
}

func (pointer GithubOauth2Client) GetConfig() *oauth2.Config{
	return pointer.config
}

func (pointer GithubOauth2Client) GetResponseURL() string{
	return pointer.responseURL
}