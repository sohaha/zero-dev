package model

import "github.com/sohaha/zlsgo/zarray"

var globalModels = zarray.NewHashMap[string, *Model]()

func GetModel(name string) (*Model, bool) {
	return globalModels.Get(name)
}
