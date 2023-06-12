package types

import (
	"github.com/gin-gonic/gin"
	"github.com/oaago/server/v2/http/middlewares/cors"
	"github.com/oaago/server/v2/http/middlewares/limiter"
	"github.com/oaago/server/v2/http/middlewares/recovery"
	"github.com/oaago/server/v2/http/middlewares/tracerid"
)

func (c *Context) Return(arg ...interface{}) {
	var code = 200
	var message interface{}
	var data interface{}
	for _, value := range arg {
		switch value.(type) {
		case string:
			message = value
		case int:
			code = value.(int)
		default:
			data = value
		}
	}
	if HttpCode[code] != nil {
		message = HttpCode[code]
	}
	// 简单处理一下 http code
	if code > 200 {
		c.JSON(code, ReturnType{
			Code:    code,
			Message: message,
			Data:    data,
		})
	} else {
		c.JSON(200, ReturnType{
			Code:    code,
			Message: message,
			Data:    data,
		})
	}
	c.Abort()
}

func (m *Middleware) AddGlobalMiddleware(fn func(ctx *Context)) {
	m.GlobalMiddleware = append(m.GlobalMiddleware, fn)
}

func (m *Middleware) AddGinGlobalMiddleware(fn func(ctx *gin.Context)) {
	m.GinGlobalMiddleware = append(m.GinGlobalMiddleware, fn)
}

func (m *Middleware) InitGinMid() {
	m.GinGlobalMiddleware = append(m.GinGlobalMiddleware, limiter.CookiesLimiter, recovery.Recovery, tracerid.TracerId, limiter.NewRateLimiterIp, limiter.NewRateLimiterUrl, cors.Cors("*"))
}
