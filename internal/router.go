package app

import (
	"zlsapp/common/restapi"
	"zlsapp/internal/account"
	"zlsapp/internal/open"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
)

// pubPath todo 后台前端界面不需要权限
var pubPath = []string{"/", "/manage/base/login", "/admin*", "/static/*", "/html/*"}

func InitRouter(_ *service.Conf) []service.Router {
	return []service.Router{
		&open.Open{
			Path: "/_",
		},
		&account.Account{
			Path: "/manage/base",
		},
		&account.Role{
			Path: "/manage/account/role",
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
		account.NewMiddleware(app, pubPath),
	}
}

func InitRouterBefore(conf *service.Conf, app *service.App) service.RouterAfter {
	return func(r *znet.Engine, app *service.App) {
		bindModelTemplate(r, app.Di)
		bindStatic(r)
	}
}
