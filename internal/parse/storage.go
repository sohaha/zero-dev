package parse

import (
	"github.com/sohaha/zlsgo/ztype"
)

type StorageType uint8

const (
	SQLStorage StorageType = iota + 1
	NoSQLStorage
)

type StorageOptionFn func(*StorageOptions) error
type StorageOptions struct {
	Fields  []string
	Limit   int
	OrderBy map[string]int8
}

type Storageer interface {
	FindOne(filter ztype.Map, fn ...StorageOptionFn) (ztype.Map, error)
	Find(filter ztype.Map, fn ...StorageOptionFn) (ztype.Maps, error)
	Migration(model *Model) Migrationer
	Insert(data ztype.Map) (lastId interface{}, err error)
	Delete(filter ztype.Map, fn ...StorageOptionFn) (int64, error)
	Update(data ztype.Map, filter ztype.Map, fn ...StorageOptionFn) (int64, error)
}

type Migrationer interface {
	Auto(deleteColumn bool) (err error)
	HasTable() bool
}
