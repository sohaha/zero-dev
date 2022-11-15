package controller

import (
	"net/http"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zarray"
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

	// 后台前端
	{
		localFileExist := zarray.NewHashMap[string, []byte]()
		fileserver := zstatic.NewFileserver("dist", func(c *znet.Context, name string, content []byte, err error) bool {
			if err != nil {
				return false
			}

			b, ok := localFileExist.ProvideGet(name, func() ([]byte, bool) {
				path := "dist/" + name
				if zfile.FileExist(path) {
					return nil, true
				}
				b, err := zfile.ReadFile(path)
				if err != nil {
					return nil, false
				}
				return b, true
			})

			if ok && content != nil {
				m := zfile.GetMimeType(name, b)
				c.Byte(200, b)
				c.SetContentType(m)
				return true
			}

			return false
		})
		r.GET(`/admin{file:.*}`, fileserver)
	}

	r.NotFoundHandler(func(c *znet.Context) {
		c.JSON(http.StatusNotFound, znet.ApiData{Code: 404, Msg: "此路不通"})
	})

}

func (h *Home) Get(c *znet.Context) {
	c.ApiJSON(200, "Success", map[string]interface{}{
		"name": h.Conf.Base.Name,
	})
}
