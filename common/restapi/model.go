package restapi

import (
	"zlsapp/internal/error_code"
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
)

type RestApi struct {
	service.App
	Path string
}

func (h *RestApi) Init(g *znet.Engine) {
	g.Use(func(c *znet.Context) error {
		name := c.GetParam("model")
		m, ok := parse.GetModel(name)
		if !ok {
			return error_code.NotFound.Text("模型不存在")
		}

		c.WithValue("model", m)
		c.Next()
		return nil
	})

	g.GET(":model", func(c *znet.Context) (interface{}, error) {
		m := c.MustValue("model").(*parse.Modeler)

		return parse.RestapiGetPage(c, m)
	})

	g.GET(":model/:key", func(c *znet.Context) (interface{}, error) {
		m := c.MustValue("model").(*parse.Modeler)

		return parse.RestapiGetInfo(c, m)
	})

	// g.POST(":model", func(c *znet.Context) (interface{}, error) {
	// 	return c.MustValue("model").(*Model).restApiCreate(c)
	// })

	// g.PATCH(":model/:key", func(c *znet.Context) (interface{}, error) {
	// 	return c.MustValue("model").(*Model).restApiUpdate(c)
	// })

	// g.DELETE(":model/:key", func(c *znet.Context) (interface{}, error) {
	// 	return c.MustValue("model").(*Model).restApiDelete(c)
	// })

}
