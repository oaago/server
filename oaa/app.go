package oaa

type AppType struct {
	Server      *OAAServer
	Router      *ConfigRouter
	ServerHooks ServerHooks
}

var App = &AppType{}
