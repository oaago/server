package socket

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/oaago/cloud/logx"
	"github.com/oaago/cloud/op"
	"github.com/oaago/server/v2/http/middlewares/cors"
	"github.com/olahol/melody"
)

func InitSocket(r *gin.Engine) {
	if op.ConfigData.Socket.Enable != true {
		return
	}
	if op.ConfigData.Socket.Types == "socketio" {
		baseUrl := op.ConfigData.Socket.BaseUrl
		socket := socketio.NewServer(nil)
		r.GET(baseUrl, cors.Cors("*"), gin.WrapH(socket))
		r.POST(baseUrl, cors.Cors("*"), gin.WrapH(socket))
	} else if op.ConfigData.Socket.Types == "ws" {
		baseUrl := op.ConfigData.Socket.BaseUrl
		m := melody.New()
		r.GET(baseUrl, cors.Cors("*"), func(c *gin.Context) {
			m.HandleRequest(c.Writer, c.Request)
		})
		r.POST(baseUrl, cors.Cors("*"), func(c *gin.Context) {
			m.HandleRequest(c.Writer, c.Request)
		})
	} else {
		logx.Logger.Info("不支持。。。")
	}
}
