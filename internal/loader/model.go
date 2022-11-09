package loader

import (
	"zlsapp/internal/model"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/zlsgo/zdb"
)

type Modeler struct {
	files
}

func (l *Loader) newModeler() {
	if l.err != nil {
		return
	}

	m := &Modeler{}
	models := make(map[string]*model.Model, 0)

	_, err := l.Di.Invoke(func(db *zdb.DB, c *service.Conf) {
		conf := c.Core()
		m.Files = Scan("./app/", Model)
		for name, path := range m.Files {
			json, err := zfile.ReadFile(path)
			if err != nil {
				l.err = zerror.With(err, "读取模型文件失败: "+path)
				return
			}
			mv, err := model.Add(db, name, json, false)
			if err != nil {
				l.err = zerror.With(err, "添加模型失败: "+path)
				return
			}
			models[path] = mv
		}

		for path, v := range models {
			err := v.Migration(conf.GetBool("migration.delete_column")).Auto()
			if err != nil {
				l.err = zerror.With(err, "模型迁移失败: "+path)
				return
			}
		}
	})

	if err != nil {
		l.err = err
	}

	l.Model = m

	return
}
