package parse

import (
	"github.com/zlsgo/zdb"
)

type SQL struct {
	db    *zdb.DB
	table string
}

// var _ storage.Storageer = (*SQL)(nil)

func NewSQL(db *zdb.DB, table string) Storageer {
	return &SQL{
		db:    db,
		table: table,
	}
}

// var _ storage.Migrationer = (*Migration)(nil)

func (s *SQL) Migration(model *Model) Migrationer {
	return &Migration{
		Model: model,
		DB:    s.db,
	}
}
