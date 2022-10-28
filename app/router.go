package app

import (
	"time"
	"zlsapp/app/model"
	"zlsapp/controller"
	"zlsapp/grbac"
	"zlsapp/grbac/meta"
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
	// grbacLoader := grbac.WithLoader(func() (grbac.Rules, error) {
	// 	rules := meta.Rules{}
	// 	zlog.Debug("重新")
	// 	return rules, nil
	// }, time.Second*10)
	grbacLoader := grbac.WithYAML("grbac/testdata/grbac.yml", time.Second*2)

	return []znet.Handler{
		cors.Default(),
		grbac.NewMiddleware(grbacLoader, grbac.WithMatchMode(meta.MatchPriorityAllow)),
	}
}
