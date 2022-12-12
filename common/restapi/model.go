package restapi

import (
	"strings"
	"zlsapp/internal/account"
	"zlsapp/internal/error_code"
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type RestApi struct {
	service.App
	Path string
}

func (h *RestApi) Init(g *znet.Engine) {
	g.Use(func(c *znet.Context) (err error) {
		name := c.GetParam("model")
		m, ok := parse.GetModel(name)
		if !ok {
			return error_code.NotFound.Text("模型不存在")
		}

		c.WithValue("model", m)
		c.Next()
		account.WrapActionLogs(c, "模型处理", m.Name)
		return nil
	})

	g.GET(":model", func(c *znet.Context) (interface{}, error) {
		m := c.MustValue("model").(*parse.Modeler)
		fields := parse.GetViewFields(m, "lists")
		var withFilds []string
		if with, _ := c.GetQuery("with"); with != "" {
			withFilds = strings.Split(with, ",")
		}
		return parse.RestapiGetPage(c, m, ztype.Map{}, fields, withFilds)
	})

	g.GET(":model/:key", func(c *znet.Context) (interface{}, error) {
		m := c.MustValue("model").(*parse.Modeler)

		fields := parse.GetViewFields(m, "info")
		var withFilds []string
		if with, _ := c.GetQuery("with"); with != "" {
			withFilds = strings.Split(with, ",")
		}
		return parse.RestapiGetInfo(c, m, fields, withFilds)
	})

	g.POST(":model", func(c *znet.Context) (interface{}, error) {
		m := c.MustValue("model").(*parse.Modeler)

		return parse.RestapiCreate(c, m)
	})

	g.PATCH(":model/:key", func(c *znet.Context) (interface{}, error) {
		m := c.MustValue("model").(*parse.Modeler)

		return parse.RestapiUpdate(c, m)
	})

	g.DELETE(":model/:key", func(c *znet.Context) (interface{}, error) {
		m := c.MustValue("model").(*parse.Modeler)

		return parse.RestapiDelete(c, m)
	})
}
