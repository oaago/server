package app

import (
	"github.com/jinzhu/copier"
	"github.com/oaago/cloud/op"
	"github.com/oaago/cloud/preload"
	"github.com/oaago/server/v2/http/core"
	"github.com/oaago/server/v2/http/event"
	"github.com/oaago/server/v2/types"
)

type Application types.Application

func (app *Application) Create() *types.Application {
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
	httpOptions := types.HttpConfig{
		Port:     op.ConfigData.Server.Port,
		EventBus: app.EventBus,
	}
	httpRouter := core.NewRouter(httpOptions)
	if app.LifeCycle.AfterLoadRouter != nil {
		app.LifeCycle.AfterLoadRouter()
	}
	app.Start = httpRouter.Start
	app.HttpEngine = (*types.HttpEngine)(httpRouter)
	app.Config = op.ConfigData
	types.App = (*types.Application)(app)
	return types.App
}
