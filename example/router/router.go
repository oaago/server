package router

import (
	"github.com/oaago/server/example/handler"
	v2 "github.com/oaago/server/v2"
)

func LoadRouter(http *v2.HttpEngine) {
	http.Router.GET("/aaa", v2.NewHandler(handler.App))
}
