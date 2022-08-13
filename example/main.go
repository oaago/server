package main

import (
	"fmt"
	"github.com/oaago/server/v2"
)

func main() {
	op := v2.HttpConfig{
		Addr: "0.0.0.0",
		Port: 8088,
	}
	http := v2.NewRouter(op)
	http.Router.GET("/aaa", v2.NewHandler(func(c *v2.Context) {
		fmt.Print("111")
	}))
	http.Start()
}
