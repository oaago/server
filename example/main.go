package main

import (
	"github.com/oaago/server/example/router"
	"github.com/oaago/server/v2/http/core"
	"github.com/oaago/server/v2/types"
)

func main() {
	op := &types.HttpConfig{
		Host: "0.0.0.0",
		Port: 8088,
	}
	http := core.NewRouter(op)
	router.LoadRouter(http)
	http.Start()
}
