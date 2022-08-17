package v2

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

type ReturnType struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func (c *Context) Return(arg ...interface{}) {
	var code = 200
	var message interface{}
	var data interface{}
	for _, value := range arg {
		switch value.(type) {
		case string:
			message = value
		case int:
			code = value.(int)
		default:
			data = value
		}
	}
	if HttpCode[code] != nil {
		message = HttpCode[code]
	}
	c.JSON(code, ReturnType{
		Code:    code,
		Message: message,
		Data:    data,
	})
	c.Abort()
}
