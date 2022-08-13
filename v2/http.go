package v2

import "github.com/gin-gonic/gin"

var HttpCode = make(map[int]interface{})

type HttpEngine struct {
	Router  *gin.Engine
	Options HttpConfig
}

type HttpConfig struct {
	GlobalMiddleware []func(ctx *Context)
	Host             string
	Port             int
	Name             string
	HttpCode         map[int]interface{}
}
