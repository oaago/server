package v2

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/oaago/server/event"
	"github.com/oaago/server/oaa/translator"
	"strconv"
)

func NewRouter(options HttpConfig) *HttpEngine {
	if options.EventBus == nil {
		options.EventBus = event.NewEvent()
	}
	options.EventBus.Publish("initRouter")
	r := gin.New()
	// 装载中间件
	for _, handlerType := range options.GlobalMiddleware {
		r.Use(NewHandler(handlerType))
	}
	options.Middleware.AddInsideMid()
	// 装载内置中间件
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
func (h *HttpEngine) Start() {
	err := translator.InitTrans("zh")
	if err != nil {
		return
	}
	for _, plugin := range h.Options.Plugins {
		plugin.Install(h)
	}
	HttpCode = h.Options.HttpCode
	go func() {
		h.Options.EventBus.Publish("startEnd")
	}()
	//e := h.Router.Run(h.Options.Host + ":" + strconv.Itoa(h.Options.Port))
	server := endless.NewServer(h.Options.Host+":"+strconv.Itoa(h.Options.Port), h.Router)
	e := server.ListenAndServe()
	if e != nil {
		go func() {
			h.Options.EventBus.Publish("startError")
		}()
		//panic(e)
	}
}
