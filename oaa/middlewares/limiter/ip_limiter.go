package limiter

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/oaago/cloud/logx"
	redis2 "github.com/oaago/cloud/redis"

	"github.com/gin-gonic/gin"
)

func NewRateLimiterIp(c *gin.Context) {
	redisClient := redis2.RedisClient.Client
	// 获取Ip白名单
	ipstr := ""
	if ipstr == "::1" {
		ipstr, _ = c.Cookie("user-agent-id")
	}
	killkey := "oaaPeriodLimit:ip:" + ipstr
	killval, err := redisClient.Get(killkey).Result()
	if len(killval) > 0 {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"status":  http.StatusTooManyRequests,
			"code":    400,
			"data":    "",
			"message": "kill all",
		})
		c.Abort()
	}
	// 针对 ip + host + 路径的限制
	key := "oaaPeriodLimit:ip:" + c.Request.Host + ":" + c.Request.Method + ":" + c.Request.URL.String()
	val, err := redisClient.Get(key).Result()
	if err != nil {
		logx.Logger.Error(err.Error())
		return
	}
	if len(val) > 0 {
		arg := strings.Split(val, ":")
		if len(arg) == 2 {
			t, _ := strconv.Atoi(arg[0])
			slidingWindow := time.Duration(t) * time.Second
			limit, _ := strconv.Atoi(arg[0])
			now := time.Now().UnixNano()
			userCntKey := fmt.Sprint(c.ClientIP(), ":", key)
			redisClient.ZRemRangeByScore(userCntKey,
				"0",
				fmt.Sprint(now-(slidingWindow.Nanoseconds()))).Result()
			reqs, _ := redisClient.ZRange(userCntKey, 0, -1).Result()
			if len(reqs) >= limit {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"status":  http.StatusTooManyRequests,
					"code":    500,
					"data":    "",
					"message": "too many request",
				})
				c.Abort()
			}
			c.Next()
			redisClient.ZAddNX(userCntKey, redis.Z{Score: float64(now), Member: float64(now)})
			redisClient.Expire(userCntKey, slidingWindow)
		}
	}
	c.Next()
}
