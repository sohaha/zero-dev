package app

import (
	"zlsapp/app/model"
	"zlsapp/controller"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
)

func InitRouter(_ *service.Conf) []service.Router {
	return []service.Router{
		&controller.Home{},
		model.NewRestApi(),
	}
}

func InitMiddleware(conf *service.Conf, app *service.App) []znet.Handler {
	return []znet.Handler{
		cors.Default(),
	}
}
