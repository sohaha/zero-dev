package model

import (
	"zlsapp/internal/error_code"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
)

type (
	Model struct {
		DB      *zdb.DB
		Name    string    `json:"name"`
		Path    string    `json:"-"`
		Table   Table     `json:"table"`
		Columns []*Column `json:"columns"`
		// Relations []*relation   `json:"relations"`
		Values        []interface{} `json:"values"`
		columnsKeys   []string
		readOnlyKeys  []string
		cryptKeys     map[string]cryptProcess
		beforeProcess map[string][]beforeProcess
		afterProcess  map[string][]afterProcess
		Options       struct {
			Api              interface{} `json:"api"`
			ApiPath          string      `json:"api_path"`
			CryptID          bool        `json:"crypt_id"`
			DisabledMigrator bool        `json:"disabled_migrator"`
			SoftDeletes      bool        `json:"soft_deletes"`
			Timestamps       bool        `json:"timestamps"`
		} `json:"options"`
	}

	Table struct {
		Name    string `json:"name"`
		Comment string `json:"comment"`
	}
)

func (m *Model) Migration(deleteColumn bool) *Migration {
	return &Migration{
		Model:  *m,
		Delete: deleteColumn,
	}
}

func (m *Model) Insert(data ztype.Map) (lastId int64, err error) {
	data, err = m.valuesBeforeProcess(data)
	if err != nil {
		return 0, err
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
	rows, err := m.DB.FindAll(m.Table.Name, func(b *builder.SelectBuilder) error {
		if !force && m.Options.SoftDeletes {
			b.Where(b.EQ(DeletedAtKey, 0))
		}
		return fn(b)
	})

	if len(m.afterProcess) > 0 {
		for i := range rows {
			row := rows[i]
			for k, v := range m.afterProcess {
				if _, ok := row[k]; ok {
					row[k], err = v[0](row.Get(k).String())
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return rows, err
}

func (m *Model) FindOne(fn func(b *builder.SelectBuilder) error, force bool) (ztype.Map, error) {
	rows, err := m.Find(func(b *builder.SelectBuilder) error {
		b.Limit(1)
		return fn(b)
	}, force)

	if err == nil && rows.Len() > 0 {
		return rows[0], nil
	}

	if err == zdb.ErrRecordNotFound {
		return ztype.Map{}, zerror.Wrap(err, zerror.ErrCode(error_code.NotFound), "")
	}
	return ztype.Map{}, err
}

func (m *Model) Update(data ztype.Map, fn func(b *builder.UpdateBuilder) error) (int64, error) {
	ndata, err := m.valuesBeforeProcess(data)
	if err != nil {
		return 0, err
	}

	return m.DB.Update(m.Table.Name, ndata, fn)
}
