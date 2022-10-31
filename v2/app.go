package v2

import (
	"github.com/jinzhu/copier"
	"github.com/oaago/cloud/op"
	"github.com/oaago/cloud/preload"
	"github.com/oaago/server/v2/event"
	"time"
)

type Application struct {
	AppId     string
	AppName   string
	Config    *op.Config
	StartTime time.Duration
	EventBus  event.Event
	LifeCycle LifeCycleType
	*HttpEngine
	RpcEngine interface{}
	Start     func()
}

type LifeCycleType struct {
	BeforeLoadConfig func()
	AfterLoadConfig  func(*op.Config)
	BeforeLoadRouter func()
	AfterLoadRouter  func()
	BeforeHttpRun    func(*op.Config)
	AfterHttpRun     func(*op.Config)
	BeforeAppStart   func()
	AfterAppStart    func()
	ExitApp          func()
}

var App *Application

func (app *Application) CreateApp() *Application {
	app.EventBus = event.NewEvent()

	if app.LifeCycle.BeforeLoadConfig != nil {
		app.LifeCycle.BeforeLoadConfig()
	}
	op.ConfigData = preload.LoadConfig()
	copier.CopyWithOption(&op.ConfigData, &app.Config, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if app.LifeCycle.AfterLoadConfig != nil {
		app.LifeCycle.AfterLoadConfig(op.ConfigData)
	}

	if app.LifeCycle.BeforeLoadRouter != nil {
		app.LifeCycle.BeforeLoadRouter()
	}
	httpOptions := HttpConfig{
		Port:     op.ConfigData.Server.Port,
		EventBus: app.EventBus,
	}
	httpRouter := NewRouter(httpOptions)
	if app.LifeCycle.AfterLoadRouter != nil {
		app.LifeCycle.AfterLoadRouter()
	}
	app.Start = httpRouter.Start
	app.HttpEngine = httpRouter
	app.Config = op.ConfigData
	App = app
	return app
}
