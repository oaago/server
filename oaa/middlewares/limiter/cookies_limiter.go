package limiter

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
)

func CookiesLimiter(c *gin.Context) {
	// 生成用户id 在这个地方可以增加规则防止
	UserAgent, _ := c.Cookie("user-agent-id")
	if len(UserAgent) == 0 {
		id, _ := uuid.GenerateUUID()
		c.SetCookie("user-agent-id", id, -1, c.Request.URL.String(), c.Request.Host, false, true)
	}
	c.Next()
}
