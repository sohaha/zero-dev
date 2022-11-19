package model

import (
	"zlsapp/internal/error_code"
	"zlsapp/internal/model/storage"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
)

type (
	Model struct {
		Schema  string `json:"$schema"`
		Raw     []byte
		DB      *zdb.DB
		Storage storage.Storageer
		Name    string           `json:"name"`
		Path    string           `json:"-"`
		Table   Table            `json:"table"`
		Columns []*Column        `json:"columns"`
		Views   map[string]*View `json:"views"`
		// views         ztype.Map
		Relations     map[string]*relation `json:"relations"`
		Values        []interface{}        `json:"values"`
		fields        []string
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

func (m *Model) Migration() (*Migration, error) {
	// s, ok := m.Storage.(*sql.SQL)
	// if !ok {
	return nil, ErrNotMigration
	// }
	// return s.NewMigration(), nil
}

func (m *Model) Insert(data ztype.Map) (lastId interface{}, err error) {
	data, err = m.valuesBeforeProcess(data)
	if err != nil {
		return 0, err
	}

	data, err = CheckData(data, m.Columns, activeCreate)
	if err != nil {
		return 0, err
	}

	if m.Options.Timestamps {
		data[CreatedAtKey] = ztime.Time()
		data[UpdatedAtKey] = ztime.Time()
	}

	if m.Options.SoftDeletes {
		data[DeletedAtKey] = 0
	}

	return m.Storage.Insert(data)
	// return m.DB.InsertMaps(m.Table.Name, data)
}

func (m *Model) Find(fn func(b *builder.SelectBuilder) error, force bool) (ztype.Maps, error) {
	rows, err := m.DB.Find(m.Table.Name, func(b *builder.SelectBuilder) error {
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
	data = filterDate(data, m.readOnlyKeys)

	data, err := m.valuesBeforeProcess(data)
	if err != nil {
		return 0, err
	}

	data, err = CheckData(data, m.Columns, activeUpdate)
	if err != nil {
		return 0, err
	}

	if m.Options.Timestamps {
		data[UpdatedAtKey] = ztime.Time()
	}

	return m.DB.Update(m.Table.Name, data, fn)
}
