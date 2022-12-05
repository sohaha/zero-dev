package loader

import (
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/zdb"
)

// import (
// 	"zlsapp/internal/mm"
// 	"zlsapp/internal/parse/storage"
// 	"zlsapp/internal/parse/storage/sql"
// 	"zlsapp/service"

// 	"github.com/sohaha/zlsgo/zerror"
// 	"github.com/sohaha/zlsgo/zfile"
// 	"github.com/sohaha/zlsgo/zlog"
// 	"github.com/zlsgo/zdb"
// )

type Modeler struct {
	files
}

func (l *Loader) loadModeler(dir ...string) *Modeler {
	m := &Modeler{}

	if l.err != nil {
		return m
	}

	models := make(map[string]*parse.Modeler, 0)

	_, err := l.Di.Invoke(func(db *zdb.DB, c *service.Conf) {
		conf := c.Core()
		path := "./app/" + Model.Dir()
		if len(dir) > 0 {
			path = dir[0]
		}
		m.Files, path = Scan(path, Model.Suffix(), true)
		l.Watch(path)

		for name, path := range m.Files {
			safePath := zfile.SafePath(path)
			json, err := zfile.ReadFile(path)
			if err != nil {
				l.err = zerror.With(err, "读取模型文件失败: "+safePath)
				return
			}
			mv, err := parse.AddModel(name, json, func(m *parse.Modeler) (parse.Storageer, error) {
				// 因为模型文件可能和内置模型重名，所以这里需要追加前缀
				m.Table.Name = "model_" + m.Table.Name
				return parse.NewSQL(db, m.Table.Name), nil
			}, false)
			if err != nil {
				l.err = zerror.With(err, "添加模型失败: "+safePath)
				return
			}

			modelLog("Register Model: " + zlog.Log.ColorTextWrap(zlog.ColorLightGreen, name))

			mv.Path = path
			models[safePath] = mv
		}

		for path, m := range models {
			err := m.Migration().Auto(conf.GetBool("migration.delete_column"))
			if err != nil {
				l.err = zerror.With(err, "模型迁移失败: "+path)
				return
			}

		}
	})

	if err != nil {
		l.err = err
	}

	if l.Model == nil {
		l.Model = m
	} else {
		for k, v := range m.Files {
			l.Model.Files[k] = v
		}
	}

	return m
}
