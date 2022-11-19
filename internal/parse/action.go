package parse

import (
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"golang.org/x/exp/constraints"
)

type Filter interface {
	ztype.Map | constraints.Integer
}

func getFilter[T Filter](filter T) ztype.Map {
	var v interface{} = filter
	if val, ok := v.(ztype.Map); ok {
		return val
	}
	return ztype.Map{
		IDKey: filter,
	}
}

func Find[T Filter](m *Model, filter T, fn ...StorageOptionFn) (ztype.Maps, error) {
	return m.Storage.Find(getFilter(filter), fn...)
}

func FindOne[T Filter](m *Model, filter T, fn ...StorageOptionFn) (ztype.Map, error) {
	return m.Storage.FindOne(getFilter(filter), fn...)
}

func Insert(m *Model, data ztype.Map) (lastId interface{}, err error) {
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

func Delete[T Filter](m *Model, filter T, fn ...StorageOptionFn) (int64, error) {
	return m.Storage.Delete(getFilter(filter), fn...)
}

func Update[T Filter](m *Model, filter T, data ztype.Map, fn ...StorageOptionFn) (int64, error) {
	return m.Storage.Update(getFilter(filter), data, fn...)
}
