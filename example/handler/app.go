package handler

import (
	"fmt"
	v2 "github.com/oaago/server/v2/http/core"
)

func App(c *v2.Context) {
	fmt.Println("v2")
}
