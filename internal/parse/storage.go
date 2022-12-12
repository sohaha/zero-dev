package parse

import (
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
)

const (
	SQLStorage StorageType = iota + 1
	NoSQLStorage
)

type StorageType uint8
type StorageJoin struct {
	Table string
	As    string
	Expr  string
}

// type StorageWhere struct {
// 	Expr string
// 	// Cond  string
// 	Field string
// 	Value interface{}
// }

type StorageOptionFn func(*StorageOptions) error
type StorageOptions struct {
	Fields  []string
	Limit   int
	OrderBy map[string]int8
	Join    []StorageJoin
	// Wheres  []StorageWhere
	// DisabledSoftDeletes bool
}

type Storageer interface {
	GetStorageType() StorageType
	Find(filter ztype.Map, fn ...StorageOptionFn) (ztype.Maps, error)
	FindOne(filter ztype.Map, fn ...StorageOptionFn) (ztype.Map, error)
	Pages(page, pagesize int, filter ztype.Map, fn ...StorageOptionFn) (ztype.Maps, PageInfo, error)
	Migration(model *Modeler) Migrationer
	Insert(data ztype.Map) (lastId interface{}, err error)
	Delete(filter ztype.Map, fn ...StorageOptionFn) (int64, error)
	Update(data ztype.Map, filter ztype.Map, fn ...StorageOptionFn) (int64, error)
}

type PageInfo struct {
	zdb.Pages
}

type Migrationer interface {
	Auto(deleteColumn bool) (err error)
	HasTable() bool
}
