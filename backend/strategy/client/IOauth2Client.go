package client

import "golang.org/x/oauth2"

type IOauth2Client interface {
	GetConfig() *oauth2.Config
	GetResponseURL() string
}