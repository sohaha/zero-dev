package parse

import (
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
)

func (m *Model) FindAll(filter ztype.Map, options StorageOptions) (ztype.Maps, error) {
	return m.Storage.FindAll(filter, options)
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
}
