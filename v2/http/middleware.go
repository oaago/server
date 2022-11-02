package http

import (
	"github.com/gin-gonic/gin"
	"github.com/oaago/server/v2/http/middlewares/limiter"
	"github.com/oaago/server/v2/http/middlewares/recovery"
	"github.com/oaago/server/v2/http/middlewares/tracerid"
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

func (m *Middleware) InitMid() {
	m.GinGlobalMiddleware = append(m.GinGlobalMiddleware, limiter.CookiesLimiter, recovery.Recovery, tracerid.TracerId, limiter.NewRateLimiterIp, limiter.NewRateLimiterUrl)
}
