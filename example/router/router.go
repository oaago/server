package router

import (
	"github.com/oaago/server/example/handler"
	core "github.com/oaago/server/v2/http/core"
)

func LoadRouter(http *core.HttpEngine) {
	http.Router.GET(
		"/aaa",
		core.NewHandler(handler.App),
	)
}
