package loader

import (
	"errors"
	"strings"
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

		m.Files = Scan(path, Model.Suffix(), true)
		l.Watch(path)

		for _, path := range m.Files {
			safePath := zfile.SafePath(path)
			models[safePath], l.err = registerModel(db, path, false)
			if l.err != nil {
				return
			}
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
		l.Model.Files = append(l.Model.Files, m.Files...)
	}

	return m
}

func registerModel(db *zdb.DB, path string, force bool) (*parse.Modeler, error) {
	path = zfile.RealPath(path)
	safePath := zfile.SafePath(path)
	json, err := zfile.ReadFile(path)
	if err != nil {
		return nil, zerror.With(err, "读取模型文件失败: "+safePath)
	}

	var root string
	for _, v := range []string{"app/models", "app/modules"} {
		p := zfile.RealPath(v)
		if strings.HasPrefix(path, p) {
			root = p
		}
	}

	name := toName(path, root)
	if force {
		m, ok := parse.GetModel(name)
		if ok && m.Path != path {
			return m, errors.New("模型名称(" + name + ")与 " + zfile.SafePath(m.Path) + " 相同")
		}
	}

	mv, err := parse.AddModelForJSON(name, json, func(m *parse.Modeler) (parse.Storageer, error) {
		// 因为模型文件可能和内置模型重名，所以这里需要追加前缀
		m.Table.Name = "model_" + m.Table.Name
		return parse.NewSQL(db, m.Table.Name), nil
	}, force)
	if err != nil {
		return nil, zerror.With(err, "添加模型失败: "+safePath)
	}

	modelLog("Register Model: " + zlog.Log.ColorTextWrap(zlog.ColorLightGreen, name))
	mv.Path = path
	return mv, nil
}
