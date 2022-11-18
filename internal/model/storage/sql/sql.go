package sql

import (
	"zlsapp/internal/model/storage"

	"github.com/zlsgo/zdb"
)

type SQL struct {
	db    *zdb.DB
	table string
}

var _ storage.Storageer = (*SQL)(nil)

func New(db *zdb.DB, table string) *SQL {
	return &SQL{
		db:    db,
		table: table,
	}
}

func (s *SQL) Migration() {

}
