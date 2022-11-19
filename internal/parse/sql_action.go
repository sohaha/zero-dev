package parse

import (
	"zlsapp/internal/error_code"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
)

func (s *SQL) Insert(data ztype.Map) (lastId interface{}, err error) {
	return s.db.InsertMaps(s.table, data)
}

func (m *SQL) FindOne(filter ztype.Map, fields []string) (ztype.Map, error) {
	rows, err := m.db.FindAll(m.table, func(b *builder.SelectBuilder) error {
		b.Limit(1)
		if len(fields) > 0 {
			b.Select(fields...)
		}
		for k := range filter {
			b.Where(b.EQ(k, filter[k]))
		}
		return nil
	})

	if err == nil && rows.Len() > 0 {
		return rows[0], nil
	}

	if err == zdb.ErrRecordNotFound {
		return ztype.Map{}, zerror.Wrap(err, zerror.ErrCode(error_code.NotFound), "")
	}
	return ztype.Map{}, err
}

func (m *SQL) FindAll(filter ztype.Map, options StorageOptions) (ztype.Maps, error) {
	fields := options.Fields
	rows, err := m.db.FindAll(m.table, func(b *builder.SelectBuilder) error {
		if len(fields) > 0 {
			b.Select(fields...)
		}
		for k := range filter {
			b.Where(b.EQ(k, filter[k]))
		}

		if len(options.OrderBy) > 0 {
			for k, v := range options.OrderBy {
				if v == -1 {
					b.OrderBy(k + " DESC")
				}
			}
		}
		return nil
	})

	if err == zdb.ErrRecordNotFound {
		return ztype.Maps{}, zerror.Wrap(err, zerror.ErrCode(error_code.NotFound), "")
	}
	return rows, nil
}
