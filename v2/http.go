package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/oaago/server/v2/event"
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
}

func (h *HttpEngine) AddPlugin(li []Plugin) {
	h.Options.Plugins = append(h.Options.Plugins, li...)
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
