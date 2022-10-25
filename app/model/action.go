package model

import (
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb/builder"
)

// ActionUpdate 更新数据
func (m *Model) ActionUpdate(key interface{}, data ztype.Map) error {
	data = filterDate(data, m.readOnlyKeys)
	data, err := CheckData(data, m.Columns, activeUpdate)
	if err != nil {
		return err
	}

	if m.Options.Timestamps {
		data[UpdatedAtKey] = ztime.Time()
	}

	_, err = m.DB.Update(m.Table.Name, data, func(b *builder.UpdateBuilder) error {
		b.Where(b.EQ(IDKey, key))
		return nil
	})
	return err
}

// ActionCreate 创建数据
func (m *Model) ActionCreate(data ztype.Map) (lastId int64, err error) {
	data, err = CheckData(data, m.Columns, activeCreate)
	zlog.Log.Debug(data, err)

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

	return m.DB.InsertMaps(m.Table.Name, data)
}
