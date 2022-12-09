package loader

import (
	// "github.com/rjeczalik/notify"

	"io/fs"
	"path/filepath"
	"strings"
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/fsnotify/fsnotify"
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/zdb"
)

// type watch struct {
// 	Dirs []string
// }

// func newWatch() *watch {
// 	w := &watch{
// 		Dirs: make([]string, 0),
// 		c:    make(chan notify.EventInfo, 1),
// 	}
// 	return w
// }

func (l *Loader) Watch(dir string) {
	dir = zfile.RealPath(dir)
	dirs := []string{dir}
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})

	for _, dir := range dirs {
		if _, ok := l.watcheDir[dir]; ok {
			continue
		}
		l.watcheDir[dir] = struct{}{}
		var err error
		err = l.watcher.Add(dir)
		zlog.Debug("add", dir, err)
	}

}

func pollEvents(di zdi.Invoker, watcher *fsnotify.Watcher) {
	for {
		event, ok := <-watcher.Events
		if !ok {
			return
		}
		for _, v := range []fsnotify.Op{fsnotify.Write, fsnotify.Create, fsnotify.Rename} {
			if event.Has(v) {
				for _, v := range []FileType{Model, Flow, View} {
					if strings.HasSuffix(event.Name, v.Suffix()) {
						reRegister(di, event.Name, v)
					}
				}
			}
		}

	}
}

func reRegister(di zdi.Invoker, file string, f FileType) {
	_, _ = di.Invoke(func(db *zdb.DB, conf *service.Conf) {
		var err error
		switch f {
		case Model:
			var m *parse.Modeler
			m, err = registerModel(db, file, true)
			if err == nil {
				err = m.Migration().Auto(conf.Core().GetBool("migration.delete_column"))
			}
		}
		if err != nil {
			zlog.Error(err)
		}
	})
}
