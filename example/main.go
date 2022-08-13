package main

import (
	"github.com/oaago/server/example/router"
	"github.com/oaago/server/v2"
)

func main() {
	op := v2.HttpConfig{
		Host: "0.0.0.0",
		Port: 8088,
	}
	http := v2.NewRouter(op)
	router.LoadRouter(http)
	http.Start()
}
