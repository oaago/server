package types

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
