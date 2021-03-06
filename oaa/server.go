package oaa

import (
	"net"
	"net/http"
	"strconv"

	"github.com/oaago/cloud/logx"
	"github.com/oaago/cloud/op"
	"github.com/oaago/gateway/discovery"
	"github.com/oaago/server/oaa/utils"
	"google.golang.org/grpc"
)

type OAAServer struct {
	Version    string `json:"version"` //服务版本
	Weight     int64  `json:"weight"`  //服务权重
	HttpServer *http.Server
	RpcServer  *grpc.Server
	EtcdServer *discovery.Register
}

func Start(ops *ConfigRouter) *OAAServer {
	var RpcPort int
	if op.ConfigData == nil {
		op.ConfigData = &op.Config{}
		op.ConfigData.Server = op.Server{
			Name: "app",
			Env:  "test",
			Port: 9909,
		}
	}
	if op.ConfigData.Server.Port == 0 {
		panic("端口错误")
	}
	srv := &OAAServer{
		HttpServer: &http.Server{
			Addr:    "0.0.0.0:" + strconv.Itoa(op.ConfigData.Server.Port),
			Handler: ops.Router,
		},
	}
	utils.AppStartPrint()
	if len(op.ConfigData.Etcd.Endpoints) > 0 && op.ConfigData.Server.Port > 0 {
		go func() {
			RpcPort = op.ConfigData.Server.Port + 1000
			listen, err := net.Listen("tcp", ":"+strconv.Itoa(RpcPort))
			if err != nil {
				logx.Logger.Error("rpc net.Listen err" + err.Error())
			}
			logx.Logger.Info("http服务端口是" + strconv.Itoa(op.ConfigData.Server.Port))
			logx.Logger.Info("rpc服务端口是" + strconv.Itoa(RpcPort))
			go func() {
				srv.EtcdServer = discovery.NewRegister(op.ConfigData.Etcd.Endpoints, logx.Logx)
				srv.EtcdServer.Register(discovery.Server{
					Name:    op.ConfigData.Name,
					Addr:    "0.0.0.0:" + strconv.Itoa(RpcPort),
					Version: op.ConfigData.Server.Version,
					Weight:  op.ConfigData.Server.Weight}, 10)
				logx.Logger.Info(op.ConfigData.Name + " 注册etcd 服务成功")
				srv.BeforeLoadHook()
				ops.MapRpcRoute().Serve(listen)
			}()
		}()
	} else {
		srv.BeforeLoadHook()
	}
	go func() {
		if err := srv.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logx.Logger.Info("listen: %s\n", err)
		}
	}()
	srv.AfterLoadHook()
	App = &AppType{
		Server: srv,
		Router: ops,
	}
	return srv
}
