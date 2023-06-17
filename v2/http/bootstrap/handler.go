package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/oaago/server/v2/types"
)

func NewHandler(f func(ctx *types.Context)) func(*gin.Context) {
	return func(context *gin.Context) {
		m := &types.Context{}
		m.Context = context
		f(m)
	}
}
