package model

import (
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb/builder"
)

func (m *Model) ActionUpdate(key interface{}, data ztype.Map) error {
	data = filterDate(data, m.readOnlyKeys)
	data, err := CheckData(data, m.Columns, activeUpdate)
	if err != nil {
		return err
	}

	if m.Options.Timestamps {
		data[UpdatedAtKey] = ztime.Time()
	}

	zlog.Debug(data[UpdatedAtKey])
	_, err = m.DB.Update(m.Table.Name, data, func(b *builder.UpdateBuilder) error {
		b.Where(b.EQ(IDKey, key))
		return nil
	})
	return err
}
