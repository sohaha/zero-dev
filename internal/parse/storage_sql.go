package parse

import (
	"github.com/zlsgo/zdb"
)

type SQL struct {
	db    *zdb.DB
	table string
}

func NewSQL(db *zdb.DB, table string) Storageer {
	return &SQL{
		db:    db,
		table: table,
	}
}

func (s *SQL) Migration(model *Modeler) Migrationer {
	return &Migration{
		Model: model,
		DB:    s.db,
	}
}
