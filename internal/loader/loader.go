package loader

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
)

type files struct {
	Files map[string]string
}

type Loader struct {
	Model   *Modeler
	Di      zdi.Invoker
	err     error
	watcher *fsnotify.Watcher
	process *process
}

func Init(di zdi.Injector) *Loader {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	l := &Loader{
		Di:      di,
		watcher: watcher,
	}

	l.newModeler()
	zerror.Panic(l.err)

	go pollEvents(watcher)
	return l
}
