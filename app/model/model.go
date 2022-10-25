package model

import (
	"github.com/sohaha/zlsgo/ztime"
	"github.com/zlsgo/zdb"
)

type (
	Model struct {
		DB      *zdb.DB
		Name    string    `json:"name"`
		Table   Table     `json:"table"`
		Columns []*Column `json:"columns"`
		// Relations []*relation   `json:"relations"`
		Values       []interface{} `json:"values"`
		columnsKeys  []string
		readOnlyKeys []string
		Options      struct {
			DisabledMigrator bool        `json:"disabled_migrator"`
			Api              interface{} `json:"api"`
			ApiPath          string      `json:"api_path"`
			SoftDeletes      bool        `json:"soft_deletes"`
			Timestamps       bool        `json:"timestamps"`
		} `json:"options"`
	}

	Table struct {
		Name    string `json:"name"`
		Comment string `json:"comment"`
	}
)

func (m *Model) Migration() *Migration {
	return &Migration{
		*m,
	}
}

func (m *Model) Insert(data map[string]interface{}) (lastId int64, err error) {
	if m.Options.Timestamps {
		now := ztime.Time()
		data[CreatedAtKey] = now
		data[UpdatedAtKey] = now
	}

	if m.Options.SoftDeletes {
		data[DeletedAtKey] = 0
	}

	lastId, err = m.DB.InsertMaps(m.Table.Name, data)

	return
}
