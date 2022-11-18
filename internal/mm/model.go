package mm

import (
	"strings"
	"zlsapp/internal/model/storage"
	"zlsapp/internal/parse"

	"github.com/sohaha/zlsgo/zarray"
)

var globalModels = zarray.NewHashMap[string, *Model]()

type Model struct {
	parse.Model
}

func Get(name string) (*Model, bool) {
	return globalModels.Get(name)
}

func Add(name string, json []byte, bindStorage func(*Model) (storage.Storageer, error), force ...bool) (*Model, error) {
	p, err := parse.ParseModel(json)
	if err != nil {
		return nil, err
	}

	m := &Model{
		Model: *p,
	}

	m.Storage, err = bindStorage(m)
	if err != nil {
		return nil, err
	}
	name = strings.TrimSuffix(name, ".model.json")
	name = strings.Replace(name, "/", "-", -1)
	if _, ok := globalModels.Get(name); ok && !(len(force) > 0 && force[0]) {
		return nil, ErrModuleAlreadyExists
	}
	globalModels.Set(name, m)

	return m, nil
}
