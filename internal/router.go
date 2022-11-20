package app

import (
	"zlsapp/common/restapi"
	"zlsapp/controller"
	"zlsapp/internal/account"
	"zlsapp/internal/open"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
)

func InitRouter(_ *service.Conf) []service.Router {
	return []service.Router{
		&controller.Home{},
		&open.Open{
			Path: "/_",
		},
		&account.Account{
			Path: "/manage/base",
		},
		&restapi.ManageRestApi{
			Path: "/manage/model",
		},
		&restapi.RestApi{
			Path: "/model",
		},
		// model.NewRestApi(),
		// model.NewManageRestApi(),
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
