package core

import (
	"github.com/gin-gonic/gin"
	"github.com/oaago/server/v2/http/event"
	"github.com/oaago/server/v2/types"
)

func NewRouter(options types.HttpConfig) *HttpEngine {
	if options.EventBus == nil {
		options.EventBus = event.NewEvent()
	}
	options.EventBus.Publish("initRouter")
	r := gin.New()
	// 装载全局中间件以及拦截器，拦截器等同于全局中间件
	options.GlobalMiddleware = append(options.GlobalMiddleware, options.Interceptor...)
	for _, handlerType := range options.GlobalMiddleware {
		r.Use(NewHandler(handlerType))
	}
	// 装载内置中间件
	InitMid(&options.Middleware)
	// 自定义中间件
	for _, f := range options.Middleware.InsideMiddType {
		r.Use(f)
	}
	// 兼容gin中间件
	for _, f := range options.Middleware.GinGlobalMiddleware {
		r.Use(f)
	}
	// 框架自动加载中间件
	for _, f := range options.Middleware.GlobalMiddleware {
		r.Use(NewHandler(f))
	}
	if len(options.Host) == 0 {
		options.Host = "0.0.0.0"
	}
	if options.Port == 0 {
		options.Port = 9901
	}
	if options.BaseUrl == "" {
		options.BaseUrl = "/"
	}
	return &HttpEngine{
		Router:  r,
		Options: options,
	}
}
