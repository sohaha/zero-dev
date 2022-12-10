package loader

import (
	// "github.com/rjeczalik/notify"

	"io/fs"
	"path/filepath"
	"strings"
	"time"
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/fsnotify/fsnotify"
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zpool"
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
		// var err error
		_ = l.watcher.Add(dir)
		// _=err
		// zlog.Debug("add", dir, err)
	}

}

func (l *Loader) pollEvents(di zdi.Invoker) {
	watcher := l.watcher
	pool := zpool.New(10)
	for {
		event, ok := <-watcher.Events
		if !ok {
			return
		}
		if event.Has(fsnotify.Remove) {
			if zfile.DirExist(event.Name) {
				_ = watcher.Remove(event.Name)
				continue
			}
		}
		for _, v := range []fsnotify.Op{fsnotify.Write, fsnotify.Create} {
			if event.Has(v) {
				if zfile.DirExist(event.Name) {
					l.Watch(event.Name)
					continue
				}
				for _, v := range []FileType{Model, Flow, View} {
					t := v
					if strings.HasSuffix(event.Name, t.Suffix()) {
						_ = pool.Do(func() {
							time.Sleep(time.Second / 2)
							reRegister(di, event.Name, t)
						})
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
