package socket

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/oaago/cloud/logx"
	"github.com/oaago/cloud/op"
	"github.com/olahol/melody"
	"net/http"
	"time"
)

func InitSocket(r *gin.Engine) {
	if op.ConfigData.Socket.Enable != true {
		return
	}
	if op.ConfigData.Socket.Types == "socketio" {
		//baseUrl := op.ConfigData.Socket.BaseUrl
		socketConfig := &engineio.Options{
			PingTimeout:  7 * time.Second,
			PingInterval: 5 * time.Second,
			Transports: []transport.Transport{
				&polling.Transport{
					Client: &http.Client{
						Timeout: time.Minute,
					},
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
				&websocket.Transport{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			},
		}
		socket := socketio.NewServer(socketConfig)
		r.GET("/socket.io/*any", gin.WrapH(socket))
		r.POST("/socket.io/*any", gin.WrapH(socket))
	} else if op.ConfigData.Socket.Types == "ws" {
		baseUrl := op.ConfigData.Socket.BaseUrl
		m := melody.New()
		r.GET(baseUrl, func(c *gin.Context) {
			m.HandleRequest(c.Writer, c.Request)
		})
		r.POST(baseUrl, func(c *gin.Context) {
			m.HandleRequest(c.Writer, c.Request)
		})
	} else {
		logx.Logger.Info("不支持。。。")
	}
}
