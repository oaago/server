package types

import (
	evbus "github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin"
	"github.com/oaago/cloud/op"
	"time"
)

type Event evbus.Bus
type BusPublisher evbus.BusPublisher
type BusSubscriber evbus.BusSubscriber
type BusController evbus.BusController

type Application struct {
	AppId     string
	AppName   string
	Config    *op.Config
	StartTime time.Duration
	EventBus  Event
	LifeCycle LifeCycleType
	*HttpEngine
	Options   *HttpConfig
	RpcEngine interface{}
	Start     func()
}

type LifeCycleType struct {
	BeforeLoadConfig func()
	AfterLoadConfig  func(*op.Config)
	BeforeLoadRouter func()
	AfterLoadRouter  func()
	BeforeHttpRun    func(*op.Config)
	AfterHttpRun     func(*op.Config)
	BeforeAppStart   func()
	AfterAppStart    func()
	ExitApp          func()
}

var App *Application

var HttpCode = make(map[int]interface{})

type HttpEngine struct {
	Router  *gin.Engine
	Options *HttpConfig
}

type Plugin interface {
	Install(*HttpEngine)
}
type Middleware struct {
	GlobalMiddleware    []func(*Context)
	PartMiddleware      []func(*Context)
	GinGlobalMiddleware []func(*gin.Context)
	InsideMiddType      []func(*gin.Context)
}
type GlobalMiddleware struct{}
type PartMiddleware struct{}
type GinGlobalMiddleware struct{}
type MiddlewareMap struct {
	GlobalMiddleware
	PartMiddleware
	GinGlobalMiddleware
}

var Middlewares Middleware

type HttpConfig struct {
	Middleware       Middleware
	GlobalMiddleware []func(ctx *Context)
	Host             string
	Port             int
	Name             string
	HttpCode         map[int]interface{}
	BaseUrl          string
	Plugins          []Plugin
	EventBus         Event
	Interceptor      []func(ctx *Context)
}

type Context struct {
	*gin.Context
}

type ReturnType struct {
	Code    interface{} `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type InsideMiddType []func(*gin.Context)

type HandlerType gin.HandlerFunc

type Handler struct {
	HandlerType gin.HandlerFunc
}
