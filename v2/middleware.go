package v2

import (
	"github.com/gin-gonic/gin"
	limiter2 "github.com/oaago/server/v2/middlewares/limiter"
	"github.com/oaago/server/v2/middlewares/recovery"
	"github.com/oaago/server/v2/middlewares/tracerid"
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
	m.GinGlobalMiddleware = append(m.GinGlobalMiddleware, limiter2.CookiesLimiter, recovery.Recovery, tracerid.TracerId, limiter2.NewRateLimiterIp, limiter2.NewRateLimiterUrl)
}
