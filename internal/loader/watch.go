package loader

import (
	// "github.com/rjeczalik/notify"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/sohaha/zlsgo/zlog"
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
	// l.c
	err := l.watcher.Add(dir)
	zlog.Debug(dir, err)

}

func pollEvents(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			zlog.Success("---:", event, ok)
			if !ok {
				return
			}
			for _, v := range []fsnotify.Op{fsnotify.Write, fsnotify.Create, fsnotify.Rename} {
				if event.Has(v) {
					log.Println("modified file:", event.Name)
				}
			}

		case err, ok := <-watcher.Errors:
			zlog.Error("error:", err)
			if !ok {
				return
			}
		}
	}
}
