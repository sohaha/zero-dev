package app

import (
	"zlsapp/controller"
	"zlsapp/internal/account"
	"zlsapp/internal/loader"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
)

// pubPath todo 后台前端界面不需要权限
var pubPath = []string{"/", "/manage/base/login", "/_*", "/static/*", "/html/*"}

func InitRouters(_ *service.Conf) []service.Router {
	r := []service.Router{
		// &controller.Home{},
		&controller.Inlay{
			Path: "/_",
		},
		// &open.Open{
		// 	Path: "/__",
		// },
		// &restapi.RestApi{
		// 	Path: "/api",
		// },

		// &restapi.RestApi{
		// 	Path:     "/model",
		// 	IsManage: true,
		// },
		// &extend.File{},

		// model.NewRestApi(),
		// model.NewManageRestApi(),
	}

	r = append(r, controller.ManageRouter()...)
	return r
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
		bindTemplate(r, app.Di)

		var l *loader.Loader
		_ = app.Di.Resolve(&l)

		bindStatic(r)
	}
}
