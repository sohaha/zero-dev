package controller

import (
	"zlsapp/internal/parse/jsonschema"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zstatic"
)

type Inlay struct {
	service.App
	Path string
}

func (h *Inlay) Init(g *znet.Engine) {
	h.registerManage(g)
}

func (h *Inlay) GetSchemaModel(c *znet.Context) {
	c.Byte(200, jsonschema.GetModelSchema())
	c.SetContentType(znet.ContentTypeJSON)
}

func (h *Inlay) registerManage(g *znet.Engine) {
	logLevel := g.Log.GetLogLevel()
	g.Log.SetLogLevel(zlog.LogNot)
	defer g.Log.SetLogLevel(logLevel)

	g.GET("", func(c *znet.Context) {
		c.Redirect(c.Request.URL.Path + "/index.html")
	})

	{
		localFileExist := zarray.NewHashMap[string, []byte]()
		fileserver := zstatic.NewFileserver("dist", func(c *znet.Context, name string, content []byte, err error) bool {
			if err != nil {
				return false
			}

			b, ok := localFileExist.ProvideGet(name, func() ([]byte, bool) {
				if content != nil {
					return content, true
				}
				path := "dist/" + name

				if !zfile.FileExist(path) {
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
		g.GET(`{file:.*}`, fileserver)
	}
}
