package controller

import (
	"net/http"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zstatic"
)

type Home struct {
	service.App
}

func (h *Home) Init(r *znet.Engine) {
	// 静态资源目录，常用于放上传的文件
	r.Static("/static/", zfile.RealPathMkdir("./resource/static"))

	fileserver := zstatic.NewFileserver("dist")
	r.GET("/admin/{file:.*}", fileserver)

	r.NotFoundHandler(func(c *znet.Context) {
		c.JSON(http.StatusNotFound, znet.ApiData{Code: 404, Msg: "此路不通"})
	})

}

func (h *Home) Get(c *znet.Context) {
	c.ApiJSON(200, "Success", map[string]interface{}{
		"name": h.Conf.Base.Name,
	})
}
