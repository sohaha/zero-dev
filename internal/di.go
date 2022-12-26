package app

import (
	"zlsapp/common/hashid"
	"zlsapp/internal/loader"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zdi"
)

func InitDI() zdi.Injector {
	di := zdi.New()

	di.Map(di, zdi.WithInterface((*zdi.Injector)(nil)))

	di.Provide(service.InitConf)
	di.Provide(service.InitApp)

	di.Provide(service.InitWeb)
	di.Provide(service.InitDB)

	di.Provide(InitMiddleware)
	di.Provide(InitRouters)
	di.Provide(InitRouterBefore)
	di.Provide(InitTasks)

	di.Provide(hashid.Init)

	di.Provide(loader.Init)

	return di
}
