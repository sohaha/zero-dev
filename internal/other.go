package app

import (
	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
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
		r.GET(`/admin{file:.*}`, fileserver)
	}
}
