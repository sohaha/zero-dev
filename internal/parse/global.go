package parse

import (
	"github.com/sohaha/zlsgo/zarray"
)

var globalModels = zarray.NewHashMap[string, *Modeler]()

func GetModel(name string) (*Modeler, bool) {
	return globalModels.Get(name)
}

func ModelsForEach(fn func(key string, m *Modeler) bool) {
	globalModels.ForEach(fn)
}

func AddModelForJSON(name string, json []byte, bindStorage func(*Modeler) (Storageer, error), force ...bool) (*Modeler, error) {
	m, err := ParseModel(json)
	if err != nil {
		return nil, err
	}

	err = AddModel(name, m, bindStorage, force...)

	return m, err
}

func AddModel(name string, m *Modeler, bindStorage func(*Modeler) (Storageer, error), force ...bool) (err error) {
	InitModel(name, m)
	m.Storage, err = bindStorage(m)
	if err != nil {
		return err
	}

	if _, ok := globalModels.Get(name); ok && !(len(force) > 0 && force[0]) {
		return ErrModuleAlreadyExists
	}
	globalModels.Set(name, m)

	return nil
}
