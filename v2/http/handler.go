package http

import (
	"github.com/gin-gonic/gin"
	"github.com/oaago/server/v2/types"
)


func NewHandler(f func(c *types.Context)) func(*gin.Context) {
	return func(context *gin.Context) {
		m := &Context{}
		m.Context = context
		f((*types.Context)(m))
	}
}
