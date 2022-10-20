package app

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"

	"github.com/sohaha/zlsgo/zutil"
)

var di zdi.Injector
var c *service.Conf

func Init() (zdi.Invoker, error) {
	err := zutil.TryCatch(func() (err error) {
		di = InitDI()

		zerror.Panic(zerror.With(di.Resolve(&c), "配置读取失败"))

		return err
	})

	return di, err
}

func Start() error {
	_, err := di.Invoke(service.RunWeb)
	if err != nil {
		err = zerror.With(err, "服务启动失败")
	} else {
		_, _ = di.Invoke(service.StopWeb)
	}

	return err
}
