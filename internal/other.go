package app

import (
	"net/http"
	"strings"

	"github.com/zlsgo/jet"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zstatic"
)

func bindStatic(r *znet.Engine) {
	// 静态资源目录，常用于放上传的文件
	r.Static("/static/", zfile.RealPath("./resource/static"))

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
}

func bindModelTemplate(r *znet.Engine) {
	dir := "app/views"
	if !zfile.DirExist(dir) {
		return
	}

	j := jet.New(r, dir, func(o *jet.Options) {
	})

	_ = j.Load()

	r.SetTemplate(j)

	var mapping ztype.Map
	if zfile.FileExist(dir + "/index.view.json") {
		f, _ := zfile.ReadFile(dir + "/index.view.json")
		_ = zjson.Unmarshal(f, &mapping)
	}

	if mapping == nil {
		if zfile.FileExist(dir + "/index.jet.html") {
			r.GET("/", func(c *znet.Context) {
				c.Template(http.StatusOK, "index", ztype.Map{})
			})
		}

		// r.Any("/html/:model/:key", func(c *znet.Context) {
		// 	zlog.Debug(c.GetAllParam())
		// })

		return
	}

	for k := range mapping {
		v := mapping.Get(k)
		k := strings.TrimLeft(k, "/")
		r.GET("/"+k, func(c *znet.Context) {
			zlog.Debug(333, j, v)
			zlog.Debug(333, j.Exists(v.String()))
			zlog.Debug(333, j.Exists("pages/"+v.String()))
			// c.Template(http.StatusOK, v.String(), ztype.Map{})
		}, znet.Recovery(func(c *znet.Context, err error) {
			zlog.Error("Recovery", err)
			c.Next()
		}))
	}
}
