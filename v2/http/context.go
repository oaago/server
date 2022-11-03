package http

import "github.com/oaago/server/v2/types"

type Context types.Context

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
	if types.HttpCode[code] != nil {
		message = types.HttpCode[code]
	}
	c.JSON(code, types.ReturnType{
		Code:    code,
		Message: message,
		Data:    data,
	})
	c.Abort()
}
