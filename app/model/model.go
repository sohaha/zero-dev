package model

import (
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
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
		cryptKeys    map[string]cryptProcess
		Options      struct {
			CryptID          bool        `json:"crypt_id"`
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

func (m *Model) Insert(data ztype.Map) (lastId int64, err error) {
	for k := range m.cryptKeys {
		if _, ok := data[k]; ok {
			data[k], err = m.cryptKeys[k](data.Get(k).String())
			if err != nil {
				return 0, err
			}
		}

	}

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

func (m *Model) Find(fn func(b *builder.SelectBuilder) error, force bool) (ztype.Maps, error) {
	return m.DB.FindAll(m.Table.Name, func(b *builder.SelectBuilder) error {
		if !force && m.Options.SoftDeletes {
			b.Where(b.EQ(DeletedAtKey, 0))
		}
		return fn(b)
	})
}

func (m *Model) FindOne(fn func(b *builder.SelectBuilder) error, force bool) (ztype.Map, error) {
	return m.DB.Find(m.Table.Name, func(b *builder.SelectBuilder) error {
		if !force && m.Options.SoftDeletes {
			b.Where(b.EQ(DeletedAtKey, 0))
		}
		return fn(b)
	})
}

func (m *Model) Update(data interface{}, fn func(b *builder.UpdateBuilder) error) (int64, error) {
	return m.DB.Update(m.Table.Name, data, fn)
}
