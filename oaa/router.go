package oaa

import (
	"strings"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

type AddRoute struct {
	Func       func(ctx *Ctx)
	Middleware []func(ctx *MidHandlerFunc)
}
type Hand []func(ctx *Ctx)
type Engine struct {
	*gin.Engine
}
type MapHttpRoute map[string]Hand

type ConfigRouter struct {
	Engine       Engine
	MapHttpRoute func() MapHttpRoute
	MapRpcRoute  func() *grpc.Server
	HttpMethod   []string
	Middleware   Middleware
	Router       *gin.Engine
	RpcServer    *grpc.Server
}

func NewRouter(options *ConfigRouter) *ConfigRouter {
	var router *gin.Engine
	if options.Router != nil {
		router = options.Router
	} else {
		router = gin.New()
		// 直接增加接口文档配置
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	if len(options.HttpMethod) == 0 {
		options.HttpMethod = []string{"GET", "POST", "DELETE", "PUT", "OPTIONS", "HEAD", "UPDATE"}
	}
	// 装载内置中间件
	options.Middleware.InitMid()
	// 加载自定义中间件
	for _, f := range options.Middleware.InsideMiddType {
		router.Use(f)
	}
	// 兼容gin中间件
	for _, f := range options.Middleware.GinGlobalMiddleware {
		router.Use(f)
	}
	// 框架自动加载中间件
	for _, f := range options.Middleware.GlobalMiddleware {
		router.Use(NewHandler(f))
	}
	for s, f := range options.MapHttpRoute() {
		for _, v := range options.HttpMethod {
			if strings.Contains(s, strings.ToLower(v)) {
				if strings.Contains(s, "@/") {
					newUrl := strings.Replace(s, strings.ToLower(v)+"@", "", 1)
					mids := strings.Split(newUrl, "|")
					router.Handle(v, mids[0], ManyHandler(f)...)
				}
			}
		}
	}
	options.Router = router
	App.Router = options
	options.Engine = Engine{
		Engine: router,
	}
	return options
}

//func Handler(f func(*Ctx)) func(*gin.Context) {
//	return func(context *gin.Context) {
//		var ctx *Ctx
//		copier.Copy(ctx, context)
//		f(ctx)
//	}
//}
