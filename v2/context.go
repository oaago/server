package v2

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

type ReturnType struct {
	Code    int
	Message interface{}
	Data    interface{}
}

func (c *Context) Return(code int, data interface{}, message interface{}) {
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
