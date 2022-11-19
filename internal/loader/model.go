package loader

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

// func (l *Loader) newModeler() {
// 	if l.err != nil {
// 		return
// 	}

// 	m := &Modeler{}
// 	models := make(map[string]*mm.Model, 0)

// 	_, err := l.Di.Invoke(func(db *zdb.DB, c *service.Conf) {
// 		// conf := c.Core()
// 		var dir string
// 		m.Files, dir = Scan("./app/", Model)
// 		l.Watch(dir)
// 		for name, path := range m.Files {
// 			safePath := zfile.SafePath(path)
// 			json, err := zfile.ReadFile(path)
// 			if err != nil {
// 				l.err = zerror.With(err, "读取模型文件失败: "+safePath)
// 				return
// 			}
// 			mv, err := mm.Add(name, json, func(m *mm.Model) (storage.Storageer, error) {
// 				return sql.New(db, m.Table.Name), nil
// 			}, false)
// 			if err != nil {
// 				l.err = zerror.With(err, "添加模型失败: "+safePath)
// 				return
// 			}

// 			modelLog("Register: " + zlog.Log.ColorTextWrap(zlog.ColorLightGreen, name))

// 			// 因为模型文件可能和内置模型重名，所以这里需要追加前缀
// 			mv.Table.Name = "model_" + mv.Table.Name
// 			mv.Path = path
// 			models[safePath] = mv
// 		}

// 		for path, v := range models {
// 			zlog.Error(path, v, "需要迁移")
// 			// err := v.Migration(conf.GetBool("migration.delete_column")).Auto()
// 			// if err != nil {
// 			// 	l.err = zerror.With(err, "模型迁移失败: "+path)
// 			// 	return
// 			// }

// 		}
// 	})

// 	if err != nil {
// 		l.err = err
// 	}

// 	l.Model = m
// }
