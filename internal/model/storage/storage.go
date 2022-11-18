package storage

import (
	"github.com/sohaha/zlsgo/ztype"
)

type Storageer interface {
	FindOne()
	Migration()
	Insert(data ztype.Map) (lastId interface{}, err error)
}
