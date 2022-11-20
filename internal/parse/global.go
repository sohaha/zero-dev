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

func AddModel(name string, json []byte, bindStorage func(*Modeler) (Storageer, error), force ...bool) (*Modeler, error) {
	m, err := ParseModel(json)
	if err != nil {
		return nil, err
	}

	m.Storage, err = bindStorage(m)
	if err != nil {
		return nil, err
	}

	// name = strings.TrimSuffix(name, ".model.json")
	// name = strings.Replace(name, "/", "-", -1)
	if _, ok := globalModels.Get(name); ok && !(len(force) > 0 && force[0]) {
		return nil, ErrModuleAlreadyExists
	}
	globalModels.Set(name, m)

	return m, nil
}
