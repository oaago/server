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
	if len(options.Host) == 0 {
		options.Host = "0.0.0.0"
	}
	if options.Port == 0 {
		options.Port = 9901
	}
	return &HttpEngine{
		Router:  r,
		Options: options,
	}
}
func (h *HttpEngine) Start() {
	HttpCode = h.Options.HttpCode
	err := h.Router.Run(h.Options.Host + ":" + strconv.Itoa(h.Options.Port))
	if err != nil {
		panic(err)
	}
}
