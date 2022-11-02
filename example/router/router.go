package router

import (
	"github.com/oaago/server/example/handler"
	"github.com/oaago/server/v2/http"
)

func LoadRouter(http *http.HttpEngine) {
	http.Router.GET("/aaa", http.NewHandler(handler.App))
}
