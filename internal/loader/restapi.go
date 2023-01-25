package loader

import (
	"strings"

	"zlsapp/core/api"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/zlsgo/zdb"
)

type HTTPer struct {
	files
}

func (l *Loader) loadRestapi(dir ...string) {
	h := &HTTPer{}
	if l.err != nil {
		return
	}

	_, l.err = l.Di.Invoke(func(db *zdb.DB, c *service.Conf, r *znet.Engine) {
		path := "./app/" + HTTP.Dir()
		if len(dir) > 0 {
			path = dir[0]
		}

		h.Files = Scan(path, HTTP.Suffix(), true)
		zlog.Debug(h.Files)
	})

	if l.HTTP == nil {
		l.HTTP = h
	} else {
		l.HTTP.Files = append(l.HTTP.Files, h.Files...)
	}
}

func registerRouter(path string, r *znet.Engine) (a *api.HTTP, err error) {
	path = zfile.RealPath(path)
	safePath := zfile.SafePath(path)
	json, err := zfile.ReadFile(path)
	if err != nil {
		return nil, zerror.With(err, "读取接口文件失败: "+safePath)
	}

	var root string
	for _, v := range []string{"app/apis", "app/modules"} {
		p := zfile.RealPath(v)
		if strings.HasPrefix(path, p) {
			root = p
		}
	}

	name := toName(path, root)
	// zlog.Debug(name, string(json))
	api.ParseRestApi(name, json)
	// os.Exit(0)
	return
}
