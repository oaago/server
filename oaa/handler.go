package oaa

import (
	"github.com/gin-gonic/gin"
)

type MidHandlerFunc struct {
	gin.HandlerFunc
}

func NewHandler(f func(c *Ctx)) func(*gin.Context) {
	return func(context *gin.Context) {
		m := &Ctx{}
		m.Context = context
		//copier.Copy(&m, &context)
		f(m)
	}
}

func ManyHandler(fl Hand) []gin.HandlerFunc {
	var midList []gin.HandlerFunc
	for _, f := range fl {
		midList = append(midList, NewHandler(f))
	}
	return midList
}
