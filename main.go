package main

import (
	"zlsapp/app"
	"zlsapp/app/migration"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
)

func main() {
	var c *service.Conf

	zcli.Name = "ZlsApp"
	zcli.Logo = `
_____                   
/  _  \  ______  ______  
/  /_\  \ \____ \ \____ \ 
/    |    \|  |_> >|  |_> >
\____|__  /|   __/ |   __/ 
	\/ |__|    |__|     `
	zcli.Version = "1.0.0"
	zcli.EnableDetach = true

	di, err := app.Init()

	if err == nil {
		_ = di.Resolve(&c)

		_ = di.Resolve(&c)

		migration.RunMigrations(di)

		var router *znet.Engine
		_ = di.Resolve(&router)
		_ = router.GET("/__", func(c *znet.Context) {
			c.String(200, "OK")
		})

		err = app.Start()
	}

	if err != nil {
		if c == nil || !c.Base.Debug {
			zcli.Error(err.Error())
		} else {
			zlog.Errorf("%+v\n", err)
		}
	}
}
