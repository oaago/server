package tracerid

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"strconv"
	"time"
)

func TracerId(c *gin.Context) {
	header := c.Request.Header
	val := header.Get("oaago-tracer-id")
	if len(val) == 0 {
		uuid, _ := uuid.GenerateUUID()
		key := c.ClientIP() + "-" + strconv.Itoa(int(time.Now().Unix())) + "-" + uuid
		c.Writer.Header().Set("oaago-tracer-id", key)
		c.Next()
	}
	c.Next()
}
