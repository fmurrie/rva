package client

import (
	"rva/helper"
	"sync"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/bitbucket"
)

var multitonBitbucketOauth2Client = make(map[interface{}]*BitbucketOauth2Client)
var onceBitbucketOauth2Client = make(map[interface{}]*sync.Once)

type BitbucketOauth2Client struct {
	config      *oauth2.Config
	responseURL string
}

func GetBitbucketOauth2Client(configurationHelper *helper.ConfigurationHelper) *BitbucketOauth2Client {
	var locker sync.Mutex
	locker.Lock()
    defer locker.Unlock()
	if onceBitbucketOauth2Client[configurationHelper] == nil {
		onceBitbucketOauth2Client[configurationHelper] = &sync.Once{}
	}
	onceBitbucketOauth2Client[configurationHelper].Do(func() {
		securityHelper := helper.GetSecurityHelper()
		multitonBitbucketOauth2Client[configurationHelper] = &BitbucketOauth2Client{
			config: &oauth2.Config{
				ClientID:     securityHelper.Decrypt(configurationHelper.GetStringValueByKey("client_id")),
				ClientSecret: securityHelper.Decrypt(configurationHelper.GetStringValueByKey("client_secret")),
				RedirectURL:  securityHelper.Decrypt(configurationHelper.GetStringValueByKey("callback_path")),
				Scopes: []string{
					"email",
				},
				Endpoint: bitbucket.Endpoint,
			},
			responseURL: securityHelper.Decrypt(configurationHelper.GetStringValueByKey("response_url")),
		}
	})
	return multitonBitbucketOauth2Client[configurationHelper]
}

func (pointer BitbucketOauth2Client) GetConfig() *oauth2.Config {
	return pointer.config
}

func (pointer BitbucketOauth2Client) GetResponseURL() string {
	return pointer.responseURL
}
