package parse

import (
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
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
	// DisabledSoftDeletes bool
}

type Storageer interface {
	Find(filter ztype.Map, fn ...StorageOptionFn) (ztype.Maps, error)
	Pages(page, pagesize int, filter ztype.Map, fn ...StorageOptionFn) (ztype.Maps, Page, error)
	Migration(model *Modeler) Migrationer
	Insert(data ztype.Map) (lastId interface{}, err error)
	Delete(filter ztype.Map, fn ...StorageOptionFn) (int64, error)
	Update(data ztype.Map, filter ztype.Map, fn ...StorageOptionFn) (int64, error)
}

type Page struct {
	zdb.Pages
}

type Migrationer interface {
	Auto(deleteColumn bool) (err error)
	HasTable() bool
}
