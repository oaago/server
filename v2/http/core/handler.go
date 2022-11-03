package core

import (
	"github.com/gin-gonic/gin"
)

func NewHandler(f func(c *Context)) func(*gin.Context) {
	return func(context *gin.Context) {
		m := &Context{}
		m.Context = context
		f(m)
	}
}
