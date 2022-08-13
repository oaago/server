package v2

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

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
	h.Router.Run(h.Options.Host + ":" + strconv.Itoa(h.Options.Port))
}
