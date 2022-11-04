package core

import (
	"github.com/oaago/server/v2/http/middlewares/cors"
	"github.com/oaago/server/v2/http/middlewares/limiter"
	"github.com/oaago/server/v2/http/middlewares/recovery"
	"github.com/oaago/server/v2/http/middlewares/tracerid"
	"github.com/oaago/server/v2/types"
)

func InitMid(m *types.Middleware) {
	m.GinGlobalMiddleware = append(m.GinGlobalMiddleware, limiter.CookiesLimiter, recovery.Recovery, tracerid.TracerId, limiter.NewRateLimiterIp, limiter.NewRateLimiterUrl, cors.Cors("*"))
}
