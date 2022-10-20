package model

import (
	"time"

	"github.com/zlsgo/zdb"
)

type (
	Model struct {
		DB      *zdb.DB
		Name    string    `json:"name"`
		Table   Table     `json:"table"`
		Columns []*Column `json:"columns"`
		// Relations []*relation   `json:"relations"`
		Values  []interface{} `json:"values"`
		Options struct {
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

	Column struct {
		Name        string        `json:"name"`
		Comment     string        `json:"comment"`
		Type        string        `json:"type"`
		Size        uint          `json:"size"`
		Tag         string        `json:"tag"`
		Nullable    bool          `json:"nullable"`
		Label       string        `json:"label"`
		Enum        interface{}   `json:"enum"`
		Default     interface{}   `json:"default"`
		Unique      interface{}   `json:"unique"`
		Index       interface{}   `json:"index"`
		Validations []validations `json:"validations"`
	}
)

func (m *Model) Migration() *Migration {
	return &Migration{
		*m,
	}
}

func (m *Model) Insert(data map[string]interface{}) (lastId int64, err error) {
	if m.Options.Timestamps {
		now := time.Now()
		data["created_at"] = now
		data["updated_at"] = now
	}

	if m.Options.SoftDeletes {
		data["deleted_at"] = 0
	}

	lastId, err = m.DB.InsertMaps(m.Table.Name, data)

	return
}
