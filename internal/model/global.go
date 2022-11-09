package model

import (
	"zlsapp/internal/error_code"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
)

var globalModels = zarray.NewHashMap[string, *Model]()

func GetModel(name string) (*Model, bool) {
	return globalModels.Get(name)
}

func modelsBindRouter(g *znet.Engine) error {
	g.Use(func(c *znet.Context) error {
		name := c.GetParam("model")
		m, ok := globalModels.Get(name)
		if !ok {
			return error_code.NotFound.Text("模型不存在")
		}

		c.WithValue("model", m)
		c.Next()
		return nil
	})

	g.GET(":model", func(c *znet.Context) error {
		return c.MustValue("model").(*Model).restApiGetPage(c)
	})

	g.GET(":model/:key", func(c *znet.Context) error {
		return c.MustValue("model").(*Model).restApiGetInfo(c)
	})

	g.POST(":model", func(c *znet.Context) error {
		return c.MustValue("model").(*Model).restApiCreate(c)
	})

	g.PATCH(":model/:key", func(c *znet.Context) error {
		return c.MustValue("model").(*Model).restApiUpdate(c)
	})

	g.DELETE(":model/:key", func(c *znet.Context) error {
		return c.MustValue("model").(*Model).restApiDelete(c)
	})

	return nil
}
