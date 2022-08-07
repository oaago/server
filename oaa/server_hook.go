package oaa

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/oaago/cloud/op"
	"github.com/oaago/server/oaa/utils"

	"github.com/oaago/cloud/logx"
	"github.com/oaago/cloud/mysql"
	"github.com/oaago/cloud/redis"
	"github.com/oaago/gateway/discovery"
)

type ServerHooks struct {
	Server     *http.Server
	EtcdServer *discovery.Register
	HandleList []func(*http.Server)
}

func NewHooks() *ServerHooks {
	return &ServerHooks{}
}

func (hooks *ServerHooks) AfterServerShutdownHook() {
	logx.Logger.Info("http服务启动成功 0.0.0.0:" + strconv.Itoa(op.ConfigData.Server.Port))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-quit
	logx.Logger.Info("即将关闭服务，请等待...")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := hooks.Server.Shutdown(ctx); err != nil {
		mysql.CloseAll()
		logx.Logger.Info("mysql服务关闭失败:")
		redis.Client.Close()
		logx.Logger.Info("redis服务关闭失败:")
		hooks.EtcdServer.Stop()
		logx.Logger.Info("etcdServer服务关闭失败:")
	} else {
		logx.Logger.Info("http服务关闭成功")
	}
	select {
	case <-ctx.Done():
		logx.Logger.Info("已经关闭")
	}
	utils.AppDownPrint()
	logx.Logger.Info("关闭成功")
}
func (hooks *ServerHooks) Append() {
	for _, f := range hooks.HandleList {
		f(hooks.Server)
	}
}

func (s *OAAServer) BeforeLoadHook() {
	hooks := reflect.TypeOf(NewHooks())
	for i := 0; i < hooks.NumMethod(); i++ {
		field := hooks.Method(i)
		if strings.Contains(strings.ToLower(field.Name), "beforeserver") {
			logx.Logger.Info(i, field.Name, field.Func)
			_ = reflect.ValueOf(&ServerHooks{}).MethodByName(field.Name).Call([]reflect.Value{})
		}
	}
}

func (s *OAAServer) AfterLoadHook() {
	hooks := reflect.TypeOf(NewHooks())
	for i := 0; i < hooks.NumMethod(); i++ {
		field := hooks.Method(i)
		if strings.Contains(strings.ToLower(field.Name), "afterserver") {
			logx.Logger.Info(i, field.Name, field.Func)
			_ = reflect.ValueOf(&ServerHooks{
				Server:     s.HttpServer,
				EtcdServer: s.EtcdServer,
			}).MethodByName(field.Name).Call([]reflect.Value{})
		}
	}
}
