package parse

import (
	"github.com/sohaha/zlsgo/zarray"
)

var globalModels = zarray.NewHashMap[string, *Model]()

func GetModel(name string) (*Model, bool) {
	return globalModels.Get(name)
}

func AddModel(name string, json []byte, bindStorage func(*Model) (Storageer, error), force ...bool) (*Model, error) {
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
