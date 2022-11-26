package controller

import (
	"html/template"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type Home struct {
	service.App
}

func (h *Home) Init(r *znet.Engine) {

	r.GET("/2", func(c *znet.Context) {
		c.Template(200, "home.html", znet.Data{"html": template.HTML("<h1>A Safe header</h1>"), "title": "ZlsGo 文档 ", "js": template.JS(`<script>console.log(` + ztype.ToString(123) + `)</script>`)})
	})
}

func (h *Home) Get(c *znet.Context) {
	c.ApiJSON(200, "Success", map[string]interface{}{
		"name": h.Conf.Base.Name,
	})
}
