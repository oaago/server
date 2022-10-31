package v2

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/oaago/cloud/logx"
	"github.com/oaago/server/v2/event"
	"github.com/oaago/server/v2/translator"
	"log"
	"strconv"
	"syscall"
)

var HttpCode = make(map[int]interface{})

type HttpEngine struct {
	Router  *gin.Engine
	Options HttpConfig
}

type Plugin interface {
	Install(*HttpEngine)
}

type HttpConfig struct {
	Middleware       Middleware
	GlobalMiddleware []func(ctx *Context)
	Host             string
	Port             int
	Name             string
	HttpCode         map[int]interface{}
	BaseUrl          string
	Plugins          []Plugin
	EventBus         event.Event
	Interceptor      []func(ctx *Context)
}

func (h *HttpEngine) AddPlugin(li []Plugin) {
	h.Options.Plugins = append(h.Options.Plugins, li...)
}

func (h *HttpEngine) AddInterceptor(li []func(ctx *Context)) {
	h.Options.Interceptor = append(h.Options.Interceptor, li...)
}

func (h *HttpEngine) SetPort(port int) {
	h.Options.Port = port
}

func (h *HttpEngine) SetMiddleware(mid Middleware) {
	h.Options.Middleware = mid
}

func (h *HttpEngine) SetBaseUrl(url string) {
	h.Options.BaseUrl = url
}

func (h *HttpEngine) AddHttpCode(codeMap map[int]interface{}) {
	if codeMap != nil {
		for i, i2 := range codeMap {
			h.Options.HttpCode[i] = i2
		}
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
	server.BeforeBegin = func(add string) {
		h.Options.EventBus.Publish("BeforeStartServer")
		if App.LifeCycle.BeforeAppStart != nil {
			App.LifeCycle.BeforeAppStart()
		}
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	e := server.ListenAndServe()
	if e != nil {
		go func() {
			h.Options.EventBus.Publish("startError")
		}()
		logx.Logger.Error(e.Error())
		//panic(e)
	}
}
