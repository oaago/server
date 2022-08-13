package v2

import "github.com/gin-gonic/gin"

type HandlerType gin.HandlerFunc

type Handler struct {
	HandlerType gin.HandlerFunc
}

func NewHandler(f func(c *Context)) func(*gin.Context) {
	return func(context *gin.Context) {
		m := &Context{}
		m.Context = context
		f(m)
	}
}
