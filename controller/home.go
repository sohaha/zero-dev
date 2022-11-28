package controller

import (
	"zlsapp/common/jet"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
)

type Home struct {
	service.App
}

func (h *Home) Init(r *znet.Engine) {

	// var views = jet.NewSet(
	// 	jet.NewOSFileSystemLoader("./resource/views"),
	// 	// jet.InDevelopmentMode(), // remove in production
	// )

	j := jet.New(r, "./common/jet/views/*.jet", ".jet")
	j.Debug(true)
	j.Reload(true).Load()
	// zlog.Debug(j.Load())
	r.SetHTMLTemplate()
	r.GET("/2", func(c *znet.Context) {
		err := j.Render(c.Writer, "index", map[string]interface{}{
			"Title": zstring.Rand(10),
		}, "layouts/main")
		c.Log.Debug(err)
		// c.Template(200, "home.html", znet.Data{"html": template.HTML("<h1>A Safe header</h1>"), "title": "ZlsGo 文档 ", "js": template.JS(`<script>console.log(` + ztype.ToString(123) + `)</script>`)})

		// view, err := views.GetTemplate("todos/index.jet")
		// if err != nil {
		// 	zlog.Println("Unexpected template err:", err.Error())
		// }
		// vars := make(jet.VarMap)
		// vars.Set("title", "title2")
		// err = view.Execute(c.Writer, vars, map[string]interface{}{
		// 	"title": "zls",
		// })
		// zlog.Debug(err)
	})
}

func (h *Home) Get(c *znet.Context) {
	c.ApiJSON(200, "Success", map[string]interface{}{
		"name": h.Conf.Base.Name,
	})
}
