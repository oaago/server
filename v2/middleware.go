package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/oaago/server/oaa/middlewares/limiter"
	"github.com/oaago/server/oaa/middlewares/recovery"
	"github.com/oaago/server/oaa/middlewares/tracerid"
)

type Middleware struct {
	GlobalMiddleware    []func(ctx *Context)
	PartMiddleware      []func(ctx *Context)
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
