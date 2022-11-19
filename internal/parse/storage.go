package parse

import (
	"github.com/sohaha/zlsgo/ztype"
)

type StorageType uint8

const (
	SQLStorage StorageType = iota + 1
	NoSQLStorage
)

type StorageOptions struct {
	Fields  []string
	OrderBy map[string]int8
}

type Storageer interface {
	FindOne(filter ztype.Map, fields []string) (ztype.Map, error)
	FindAll(filter ztype.Map, options StorageOptions) (ztype.Maps, error)
	Migration(model *Model, deleteColumn bool) Migrationer
	Insert(data ztype.Map) (lastId interface{}, err error)
}

type Migrationer interface {
	Auto() (err error)
}
