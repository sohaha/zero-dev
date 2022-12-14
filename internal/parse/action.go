package parse

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"golang.org/x/exp/constraints"
)

type Filter interface {
	ztype.Map | constraints.Integer | string
}

func getFilter[T Filter](m *Modeler, filter T) ztype.Map {
	var v interface{} = filter

	val, ok := v.(ztype.Map)
	if !ok {
		val = ztype.Map{
			m.Table.Name + "." + IDKey: filter,
		}
	}

	if m.Options.SoftDeletes {
		val[m.Table.Name+"."+DeletedAtKey] = 0
	}

	return val
}

func Pages[T Filter](m *Modeler, page, pagesize int, filter T, fn ...StorageOptionFn) (ztype.Maps, PageInfo, error) {
	rows, pages, err := m.Storage.Pages(page, pagesize, getFilter(m, filter), func(so *StorageOptions) error {
		if len(fn) > 0 {
			if err := fn[0](so); err != nil {
				return err
			}
		}

		if len(so.Fields) > 0 && len(so.Join) == 0 {
			fullFields := m.fullFields
			so.Fields = zarray.Filter(so.Fields, func(_ int, f string) bool {
				return zarray.Contains(fullFields, f)
			})
		}

		return nil
	})
	if err != nil {
		return rows, pages, err
	}

	afterProcess := m.afterProcess
	if len(afterProcess) == 0 {
		return rows, pages, nil
	}

	for i := range rows {
		row := &rows[i]
		for k, v := range afterProcess {
			if _, ok := (*row)[k]; ok {
				(*row)[k], err = v[0](row.Get(k).String())
				if err != nil {
					return rows, pages, err
				}
			}
		}
	}
	return rows, pages, nil
}

func find(m *Modeler, filter ztype.Map, fn ...StorageOptionFn) (ztype.Maps, error) {
	rows, err := m.Storage.Find(filter, func(so *StorageOptions) error {
		for i := range fn {
			if err := fn[i](so); err != nil {
				return err
			}
		}

		if len(so.Fields) > 0 && len(so.Join) == 0 {
			so.Fields = zarray.Filter(so.Fields, func(_ int, f string) bool {
				if strings.Contains(zstring.TrimSpace(f), " ") {
					return true
				}
				return zarray.Contains(m.fullFields, f)
			})
		}
		return nil
	})
	if err != nil {
		return rows, err
	}

	if len(m.afterProcess) > 0 {
		for i := range rows {
			row := &rows[i]
			for k, v := range m.afterProcess {
				if _, ok := (*row)[k]; ok {
					(*row)[k], err = v[0](row.Get(k).String())
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return rows, nil
}

func Find[T Filter](m *Modeler, filter T, fn ...StorageOptionFn) (ztype.Maps, error) {
	return find(m, getFilter(m, filter), fn...)
}

func FindOne[T Filter](m *Modeler, filter T, fn ...StorageOptionFn) (ztype.Map, error) {
	rows, err := find(m, getFilter(m, filter), func(so *StorageOptions) error {
		for i := range fn {
			if err := fn[i](so); err != nil {
				return err
			}
		}

		so.Limit = 1
		return nil
	})
	if err != nil {
		return ztype.Map{}, err
	}

	return rows.Index(0), nil
}

func Insert(m *Modeler, data ztype.Map, createdBy string) (lastId interface{}, err error) {
	data, err = m.valuesBeforeProcess(data)
	if err != nil {
		return 0, err
	}

	data, err = VerifiData(data, m.Columns, activeCreate)
	if err != nil {
		return 0, err
	}

	if m.Options.Timestamps {
		data[CreatedAtKey] = ztime.Time()
		data[UpdatedAtKey] = ztime.Time()
	}

	if m.Options.CreatedBy {
		data[CreatedByKey] = createdBy
	}

	if m.Options.SoftDeletes {
		data[DeletedAtKey] = 0
	}

	return m.Storage.Insert(data)
}

func Delete[T Filter](m *Modeler, filter T, fn ...StorageOptionFn) (int64, error) {
	if m.Options.SoftDeletes {
		return m.Storage.Update(ztype.Map{
			DeletedAtKey: ztime.Time().Unix(),
		}, getFilter(m, filter))
	}

	return m.Storage.Delete(getFilter(m, filter), fn...)
}

func Update[T Filter](m *Modeler, filter T, data ztype.Map, fn ...StorageOptionFn) (total int64, err error) {
	data = filterDate(data, m.readOnlyKeys)

	data, err = m.valuesBeforeProcess(data)

	if err != nil {
		return 0, err
	}

	data, err = VerifiData(data, m.Columns, activeUpdate)
	if err != nil {
		return 0, err
	}

	if m.Options.Timestamps {
		data[UpdatedAtKey] = ztime.Time()
	}

	return m.Storage.Update(data, getFilter(m, filter), fn...)
}
