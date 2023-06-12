package core

import (
	"github.com/fvbock/endless"
	"github.com/oaago/cloud/logx"
	"github.com/oaago/server/v2/http/translator"
	"github.com/oaago/server/v2/socket"
	"github.com/oaago/server/v2/types"
	"log"
	"strconv"
	"syscall"
)

type HttpEngine types.HttpEngine

func (h *HttpEngine) AddPlugin(li []types.Plugin) {
	h.Options.Plugins = append(h.Options.Plugins, li...)
}

func (h *HttpEngine) AddInterceptor(li []func(ctx *types.Context)) {
	h.Options.Interceptor = append(h.Options.Interceptor, li...)
}

func (h *HttpEngine) SetPort(port int) {
	h.Options.Port = port
}

func (h *HttpEngine) SetMiddleware(mid types.Middleware) {
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
		plugin.Install((*types.HttpEngine)(h))
	}
	types.HttpCode = h.Options.HttpCode
	go func() {
		h.Options.EventBus.Publish("startEnd")
	}()
	socket.InitSocket(h.Router)
	server := endless.NewServer(h.Options.Host+":"+strconv.Itoa(h.Options.Port), h.Router)
	server.BeforeBegin = func(add string) {
		h.Options.EventBus.Publish("BeforeStartServer")
		if types.App != nil && types.App.LifeCycle.BeforeAppStart != nil {
			types.App.LifeCycle.BeforeAppStart()
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
