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
	loader := &Loader{
		Di:        di,
		watcheDir: make(map[string]struct{}),
	}

	_, _ = di.Invoke(func(conf *service.Conf) {
		if conf.Base.Debug || conf.Base.Watch {
			watcher, err := fsnotify.NewWatcher()
			if err == nil {
				loader.watcher = watcher

				go loader.pollEvents(di)
			}
		}
	})

	// loader.loadViews()
	loader.loadModeler()
	loader.loadModules()
	loader.loadRestapi()

	// zlog.Panic(loader.err)
	zerror.Panic(loader.err)

	return loader
}
