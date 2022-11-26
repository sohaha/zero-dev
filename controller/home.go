package controller

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
)

type Home struct {
	service.App
}

func (h *Home) Init(r *znet.Engine) {

	r.GET("/2", func(c *znet.Context) {
		c.Template(200, "home.html", znet.Data{"title": "ZlsGo 文档 "})
	})

}

func (h *Home) Get(c *znet.Context) {
	c.ApiJSON(200, "Success", map[string]interface{}{
		"name": h.Conf.Base.Name,
	})
}
