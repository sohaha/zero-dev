package main

import (
	app "zlsapp/internal"
	"zlsapp/internal/account"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zutil"
)

var name = "ZeroApp"
var description = "ZeroApp is a web application framework."

func main() {
	var c *service.Conf

	zcli.Name = name
	zcli.Version = "1.0.0"
	zcli.EnableDetach = true

	di, err := app.Init()

	if err == nil {
		_, _ = di.Invoke(func(router *znet.Engine) {

			// if c.Base.Debug {
			// 	_ = router.GET(`/debug/statsviz{*:[\S]*}`, func(c *znet.Context) {
			// 		q := c.GetParam("*")
			// 		if q == "" {
			// 			c.Redirect("/debug/statsviz/")
			// 			return
			// 		}
			// 		if q == "/ws" {
			// 			statsviz.Ws(c.Writer, c.Request)
			// 			return
			// 		}
			// 		statsviz.IndexAtRoot("/debug/statsviz").ServeHTTP(c.Writer, c.Request)
			// 	}, znet.WrapFirstMiddleware(func(c *znet.Context) {
			// 		c.WithValue(account.DisabledAuthKey, true)
			// 		c.Next()
			// 	}))
			// }

		})

		err = zutil.TryCatch(func() error {
			_ = di.Resolve(&c)

			zcli.Add("passwd", "Modify account password", &account.PasswdCommand{
				DI:   di,
				Conf: c,
			})

			return zcli.LaunchServiceRun(zcli.Name, description, func() {
				zerror.Panic(app.Start())
			})
		})
	}
	if err != nil {
		if c == nil || !c.Base.Debug {
			zcli.Error(err.Error())
		} else {
			zlog.Errorf("%+v\n", err)
		}
	}

}
