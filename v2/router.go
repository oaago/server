package v2

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type HttpEngine struct {
	Router  *gin.Engine
	Options HttpConfig
}

type HttpConfig struct {
	GlobalMiddleware []func(ctx *Context)
	Addr             string
	Port             int
	Name             string
}

func NewRouter(options HttpConfig) *HttpEngine {
	r := gin.New()
	// 装载中间件
	for _, handlerType := range options.GlobalMiddleware {
		r.Use(NewHandler(handlerType))
	}
	return &HttpEngine{
		Router:  r,
		Options: options,
	}
}
func (h *HttpEngine) Start() {
	h.Router.Run(h.Options.Addr + ":" + strconv.Itoa(h.Options.Port))
}
