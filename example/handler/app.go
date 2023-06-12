package handler

import (
	"fmt"
	"github.com/oaago/server/v2/types"
)

type T struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func App(c *types.Context) {
	fmt.Println("v2")
	c.Return(21)
}
