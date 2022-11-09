package loader

import (
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
)

type files struct {
	Files map[string]string
}

type Loader struct {
	Model *Modeler
	Di    zdi.Invoker
	err   error
}

func Init(di zdi.Injector) *Loader {
	l := &Loader{
		Di: di,
	}
	l.newModeler()
	zerror.Panic(l.err)
	return l
}
