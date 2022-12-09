package loader

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/zdb"
)

type Views struct {
	files
}

type view struct {
}

func (l *Loader) loadViews(dir ...string) *Views {
	m := &Views{}

	if l.err != nil {
		return m
	}

	// models := make(map[string]*parse.Modeler, 0)

	_, err := l.Di.Invoke(func(db *zdb.DB, c *service.Conf) {
		// conf := c.Core()
		root := zfile.RealPath("./app/" + View.Dir())
		if len(dir) > 0 {
			root = dir[0]
		}
		m.Files = Scan(root, View.Suffix(), true)

		for _, path := range m.Files {
			name := toName(path, root)
			safePath := zfile.SafePath(path)
			json, err := zfile.ReadFile(path)
			if err != nil {
				l.err = zerror.With(err, "读取视图配置文件失败: "+safePath)
				return
			}
			_ = json
			// mv, err := parse.AddModel(name, json, func(m *parse.Modeler) (parse.Storageer, error) {
			// 	// 因为模型文件可能和内置模型重名，所以这里需要追加前缀
			// 	m.Table.Name = "model_" + m.Table.Name
			// 	return parse.NewSQL(db, m.Table.Name), nil
			// }, false)
			// if err != nil {
			// 	l.err = zerror.With(err, "添加模型失败: "+safePath)
			// 	return
			// }

			modelLog("Register View: " + zlog.Log.ColorTextWrap(zlog.ColorLightGreen, name))

			// models[safePath] = mv
		}

	})

	if err != nil {
		l.err = err
	}

	if l.Views == nil {
		l.Views = m
	} else {
		for k, v := range m.Files {
			l.Views.Files[k] = v
		}
	}

	return m
}
