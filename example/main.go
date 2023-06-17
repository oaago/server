package main

import (
	"github.com/oaago/server/example/router"
	"github.com/oaago/server/v2/app"
	"github.com/oaago/server/v2/http/bootstrap"
)

func main() {
	op := &app.HttpConfig{
		Host: "0.0.0.0",
		Port: 8088,
		HttpCode: map[int]interface{}{
			21: "shishi",
			22: "shishijiushishi",
		},
	}
	http := bootstrap.NewRouter(op)
	router.LoadRouter(http)
	http.Start()
}
