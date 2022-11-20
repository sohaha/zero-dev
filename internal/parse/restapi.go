package parse

import (
	"errors"
	"strings"
	"zlsapp/internal/error_code"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zvalid"
)

func GetPages(c *znet.Context) (page, pagesize int, err error) {
	rule := c.ValidRule().IsNumber().MinInt(1)
	err = zvalid.Batch(
		zvalid.BatchVar(&page, c.Valid(rule, "page", "页码").Customize(func(rawValue string, err error) (string, error) {
			if err != nil || rawValue == "" {
				return "1", err
			}
			return rawValue, nil
		})),
		zvalid.BatchVar(&pagesize, c.Valid(rule, "pagesize", "数量").MaxInt(1000).Customize(func(rawValue string, err error) (string, error) {
			if err != nil || rawValue == "" {
				return "10", err
			}
			return rawValue, nil
		})),
	)
	return
}

func restApiInfo(m *Modeler, key string, hasPrefix bool, fn ...StorageOptionFn) (ztype.Map, error) {
	filter := ztype.Map{}
	if key != "" && key != "0" {
		if hasPrefix {
			filter[m.Table.Name+"."+IDKey] = key
		} else {
			filter[IDKey] = key
		}
	}

	return FindOne(m, filter, fn...)
}

func RestapiGetInfo(c *znet.Context, m *Modeler) (interface{}, error) {
	key := c.GetParam("key")

	fields := GetViewFields(m, "info")
	finalFields, tmpFields, quote, with, withMany := getFinalFields(m, c, fields)

	info, err := restApiInfo(m, key, quote, func(so *StorageOptions) error {
		table := m.Table.Name
		for k, v := range with {
			m, ok := GetModel(v.Model)
			if !ok {
				return errors.New("关联模型(" + v.Model + ")不存在")
			}

			t := m.Table.Name
			asName := k
			so.Join = append(so.Join, StorageJoin{
				Table: t,
				As:    asName,
				Expr:  asName + "." + v.Foreign + " = " + table + "." + v.Key,
			})

			if len(v.Fields) > 0 {
				finalFields = append(finalFields, zarray.Map(v.Fields, func(_ int, v string) string {
					return asName + "." + v
				})...)
			} else {
				finalFields = append(finalFields, asName+".*")
			}
		}
		so.Fields = finalFields
		return nil
	})

	if err != nil {
		return nil, err
	}

	for k, v := range withMany {
		m, ok := GetModel(v.Model)
		if !ok {
			return nil, errors.New("关联模型(" + v.Model + ")不存在")
		}
		key := info.Get(v.Key)
		if !key.Exists() {
			return nil, errors.New("字段(" + v.Key + ")不存在，无法关联模型(" + v.Model + ")")
		}

		rows, _ := Find(m, ztype.Map{
			v.Foreign: key.Value(),
		}, func(so *StorageOptions) error {
			if len(v.Fields) > 0 {
				so.Fields = v.Fields
			}
			return nil
		})

		_ = info.Set(k, rows)
	}

	for _, v := range tmpFields {
		s := strings.SplitN(v, ".", 2)
		if len(s) == 2 {
			_ = info.Delete(s[1])
		} else {
			_ = info.Delete(v)
		}
	}

	return info, nil

}

func RestapiGetPage(c *znet.Context, m *Modeler) (interface{}, error) {
	page, pagesize, err := GetPages(c)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	fields := GetViewFields(m, "lists")
	finalFields, tmpFields, quote, with, withMany := getFinalFields(m, c, fields)

	filter := ztype.Map{}

	rows, pageInfo, err := Pages(m, page, pagesize, filter, func(so *StorageOptions) error {
		so.OrderBy = map[string]int8{m.Table.Name + "." + IDKey: -1}
		if quote {
			table := m.Table.Name
			for k, v := range with {
				m, ok := GetModel(v.Model)
				if !ok {
					return errors.New("关联模型(" + v.Model + ")不存在")
				}

				t := m.Table.Name
				asName := k
				so.Join = append(so.Join, StorageJoin{
					Table: t,
					As:    k,
					Expr:  asName + "." + v.Foreign + " = " + table + "." + v.Key,
				})

				if len(v.Fields) > 0 {
					finalFields = append(finalFields, zarray.Map(v.Fields, func(_ int, v string) string {
						return asName + "." + v
					})...)
				} else {
					finalFields = append(finalFields, asName+".*")
				}
			}
		}

		so.Fields = finalFields
		return nil
	})

	if err != nil {
		return nil, err
	}

	_ = withMany
	for _, info := range rows {
		for _, v := range tmpFields {
			_ = info.Delete(v)
		}
	}

	return ztype.Map{
		"items": rows,
		"page":  pageInfo,
	}, nil

}

func RestapiCreate(c *znet.Context, m *Modeler) (interface{}, error) {
	json, err := c.GetJSONs()
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	json = json.MatchKeys(m.fields)
	data := json.MapString()

	id, err := Insert(m, data)

	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	return ztype.Map{"id": id}, nil
}

func RestapiDelete(c *znet.Context, m *Modeler) (interface{}, error) {
	key := c.GetParam("key")
	_, err := restApiInfo(m, key, false)
	if err != nil {
		return nil, err
	}

	_, err = Delete(m, key, func(so *StorageOptions) error {
		so.Limit = 1
		return nil
	})

	return nil, err
}

func RestapiUpdate(c *znet.Context, m *Modeler) (interface{}, error) {
	key := c.GetParam("key")
	_, err := restApiInfo(m, key, false)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	json, err := c.GetJSONs()
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}
	json = json.MatchKeys(m.fields)

	data := json.MapString()
	if len(data) == 0 {
		return nil, error_code.InvalidInput.Text("没有可更新数据")
	}

	_, err = Update(m, key, data)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	return nil, nil
}
