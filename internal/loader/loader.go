package loader

import (
	"zlsapp/service"

	"github.com/fsnotify/fsnotify"
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
)

type files struct {
	Files []string
}

type Loader struct {
	Model     *Modeler
	Views     *Views
	Di        zdi.Invoker
	err       error
	watcher   *fsnotify.Watcher
	watcheDir map[string]struct{}
	process   *process
}

func Init(di zdi.Injector) *Loader {

	l := &Loader{
		Di:        di,
		watcheDir: make(map[string]struct{}),
	}

	_, _ = di.Invoke(func(conf *service.Conf) {
		if conf.Base.Debug || conf.Base.Watch {
			watcher, err := fsnotify.NewWatcher()
			if err == nil {
				l.watcher = watcher

				go pollEvents(di, watcher)
			}
		}
	})

	l.loadViews()
	l.loadModeler()
	l.loadModules()
	zerror.Panic(l.err)

	return l
}
