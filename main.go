package main

import (
	app "zlsapp/internal"
	"zlsapp/internal/loader"
	"zlsapp/service"

	"github.com/arl/statsviz"
	"github.com/sohaha/zlsgo/zcli"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zutil"
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

	err := zutil.TryCatch(func() error {
		di, err := app.Init()
		if err == nil {
			_, _ = di.Invoke(func(app *service.App, l *loader.Loader) {
				c = app.Conf
			})

			var router *znet.Engine
			_ = di.Resolve(&router)
			if c.Base.Debug {
				_ = router.GET(`/debug/statsviz{*:[\S]*}`, func(c *znet.Context) {
					if c.GetParam("*") == "/ws" {
						statsviz.Ws(c.Writer, c.Request)
						return
					}
					statsviz.IndexAtRoot("/debug/statsviz").ServeHTTP(c.Writer, c.Request)
				})
			}

			err = app.Start()
		}
		return err
	})

	if err != nil {
		if c == nil || !c.Base.Debug {
			zcli.Error(err.Error())
		} else {
			zlog.Errorf("%+v\n", err)
		}
	}
}
