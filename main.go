package main

import (
	"net/http"
	app "zlsapp/internal"
	"zlsapp/internal/account"
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
			_ = di.Resolve(&c)

			var router *znet.Engine
			_, _ = di.Invoke(func(r *znet.Engine, l *loader.Loader) {
				router = r
			})

			if c.Base.Debug {
				_ = router.HandleWrapBefore(http.MethodGet, `/debug/statsviz{*:.*}`, func(c *znet.Context) {
					q := c.GetParam("*")
					if q == "" {
						c.Redirect("/debug/statsviz/")
						return
					}
					if q == "/ws" {
						statsviz.Ws(c.Writer, c.Request)
						return
					}
					statsviz.IndexAtRoot("/debug/statsviz").ServeHTTP(c.Writer, c.Request)
				}, func(c *znet.Context) {
					c.WithValue(account.DisabledAuthKey, true)
					c.Next()
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
