package cors

import (
	"github.com/gin-gonic/gin"
	"github.com/oaago/cloud/op"
	"net/http"
)

func Cors(allowOrigin string) gin.HandlerFunc {
	return func(context *gin.Context) {
		Origin := context.Writer.Header().Get("Access-Control-Allow-Origin")
		if Origin == "*" || len(op.ConfigData.Server.Cors) == 0 {
			return
		} else {
			for _, v := range op.ConfigData.Server.Cors {
				if context.Request.Host == v {
					allowOrigin = context.Request.Host
					continue
				}
			}
			method := context.Request.Method
			context.Header("Access-Control-Allow-Origin", allowOrigin)
			context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, UPDATE, DELETE")
			context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			context.Header("Access-Control-Allow-Credentials", "true")
			if method == "OPTIONS" {
				context.AbortWithStatus(http.StatusNoContent)
			}
		}
		context.Next()
	}
}
