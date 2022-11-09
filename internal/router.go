package app

import (
	"zlsapp/controller"
	"zlsapp/internal/account"
	"zlsapp/internal/model"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
)

func InitRouter(_ *service.Conf) []service.Router {
	return []service.Router{
		&controller.Home{},
		&account.Account{
			Path: "/base",
		},
		model.NewRestApi(),
		model.NewManageRestApi(),
	}
}

func InitMiddleware(conf *service.Conf, app *service.App) []znet.Handler {
	// grbacLoader := grbac.WithLoader(func() (grbac.Rules, error) {
	// 	rules := meta.Rules{}
	// 	zlog.Debug("重新")
	// 	return rules, nil
	// }, time.Second*10)

	return []znet.Handler{
		cors.Default(),
		account.NewMiddleware(app),
	}
}
