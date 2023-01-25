package controller

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
)

type Home struct {
	service.App
}

func (f *Home) Init(g *znet.Engine) {
	g.Any("/*", func() (interface{}, error) {
		return "Hello World", nil
	}, znet.WrapFirstMiddleware(func(c *znet.Context) {
		// zlog.Debug(456)
		// c.WithValue(conf.DisabledAuthKey, true)
		c.Next()
	}))
}
