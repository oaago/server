package oaa

import (
	"github.com/gin-gonic/gin"
	"github.com/oaago/server/oaa/middlewares/limiter"
	"github.com/oaago/server/oaa/middlewares/recovery"
	"github.com/oaago/server/oaa/middlewares/tracerid"
	"reflect"
)

type Middleware struct {
	GlobalMiddleware    []func(*Ctx)
	PartMiddleware      []func(*Ctx)
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

type InsideMiddType []func(ctx *gin.Context)

func (m *Middleware) AddInsideMid() {
	m.InsideMiddType = []func(ctx *gin.Context){
		limiter.CookiesLimiter,    // cookies
		recovery.Recovery,         // 错误处理
		tracerid.TracerId,         // 生成唯一id
		limiter.NewRateLimiterIp,  // 用户标识限流
		limiter.NewRateLimiterUrl, // url 限流
	}
}

func MiddlewareUse(globalMiddleware GlobalMiddleware, partMiddleware PartMiddleware, ginGlobalMiddleware GinGlobalMiddleware) Middleware {
	midGlobal := reflect.TypeOf(globalMiddleware)
	midPart := reflect.TypeOf(partMiddleware)
	ginMid := reflect.TypeOf(ginGlobalMiddleware)

	for i := 0; i < midPart.NumMethod(); i++ {
		values := reflect.ValueOf(midPart).MethodByName(midPart.Method(i).Name)
		Middlewares.PartMiddleware = append(Middlewares.PartMiddleware, func(ctx *Ctx) {
			params := make([]reflect.Value, 1)
			params[0] = reflect.ValueOf(ctx)
			values.Call(params)
		})
	}
	for i := 0; i < midGlobal.NumMethod(); i++ {
		Middlewares.GlobalMiddleware = append(Middlewares.GlobalMiddleware, func(ctx *Ctx) {
			values := reflect.ValueOf(midGlobal).MethodByName(midGlobal.Method(i).Name)
			params := make([]reflect.Value, 1)
			params[0] = reflect.ValueOf(ctx)
			values.Call(params)
		})
	}
	for i := 0; i < ginMid.NumMethod(); i++ {
		Middlewares.GinGlobalMiddleware = append(Middlewares.GinGlobalMiddleware, func(ctx *gin.Context) {
			values := reflect.ValueOf(ginMid).MethodByName(ginMid.Method(i).Name)
			params := make([]reflect.Value, 1)
			params[0] = reflect.ValueOf(ctx)
			values.Call(params)
		})
	}
	return Middlewares
}
