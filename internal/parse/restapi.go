package parse

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

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

// func restApiDelete(m *Modeler,c *znet.Context) (interface{}, error) {
// 	key := c.GetParam("key")
// 	_, err := m.restApiInfo(key, false)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if m.Options.SoftDeletes {
// 		_, err = m.DB.Update(m.Table.Name, map[string]interface{}{
// 			DeletedAtKey: ztime.Time().Unix(),
// 		}, func(b *builder.UpdateBuilder) error {
// 			b.Where(b.EQ(IDKey, key))
// 			return nil
// 		})
// 	} else {
// 		_, err = m.DB.Delete(m.Table.Name, func(b *builder.DeleteBuilder) error {
// 			b.Where(b.EQ(IDKey, key))
// 			return nil
// 		})
// 	}

// 	return nil, err
// }

// func restApiGetPage(m *Modeler,c *znet.Context) (interface{}, error) {
// 	page, pagesize, err := GetPages(c)
// 	if err != nil {
// 		return nil, error_code.InvalidInput.Error(err)
// 	}

// 	fields := GetViewFields(m, "lists")
// 	finalFields, tmpFields, quote, with, withMany := getFinalFields(m, c, fields)

// 	rows, pages, err := m.DB.Pages(m.Table.Name, page, pagesize, func(b *builder.SelectBuilder) error {
// 		deletedAtKey := DeletedAtKey
// 		idKey := IDKey
// 		if quote {
// 			table := m.Table.Name
// 			for k, v := range with {
// 				m, ok := Get(v.Model)
// 				if !ok {
// 					return errors.New("关联模型(" + v.Model + ")不存在")
// 				}

// 				t := m.Table.Name
// 				asName := k
// 				b.JoinWithOption("", b.As(t, asName),
// 					asName+"."+v.Foreign+" = "+table+"."+v.Key,
// 				)
// 				if len(v.Fields) > 0 {
// 					finalFields = append(finalFields, zarray.Map(v.Fields, func(_ int, v string) string {
// 						return asName + "." + v
// 					})...)
// 				} else {
// 					finalFields = append(finalFields, asName+".*")
// 				}
// 			}
// 			b.Desc(table + "." + IDKey)
// 			deletedAtKey = m.Table.Name + "." + DeletedAtKey
// 			idKey = m.Table.Name + "." + IDKey
// 		}

// 		b.Select(finalFields...)
// 		b.Desc(idKey)
// 		if m.Options.SoftDeletes {
// 			b.Where(b.EQ(deletedAtKey, 0))
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	_ = withMany
// 	for _, info := range rows {
// 		for _, v := range tmpFields {
// 			_ = info.Delete(v)
// 		}
// 	}

// 	return ResultPages(rows, pages), nil

// }

// func restApiCreate(c *znet.Context,m *Modeler) (interface{}, error) {
// 	json, err := c.GetJSONs()
// 	if err != nil {
// 		return nil, error_code.InvalidInput.Error(err)
// 	}

// 	json = json.MatchKeys(m.fields)
// 	data := json.MapString()

// 	id, err := m.ActionCreate(data)
// 	if err != nil {
// 		return nil, error_code.InvalidInput.Error(err)
// 	}

// 	return ztype.Map{"id": id}, nil
// }

// func restApiUpdate(c *znet.Context,m *Modeler) (interface{}, error) {
// 	key := c.GetParam("key")
// 	_, err := m.restApiInfo(key, false)
// 	if err != nil {
// 		return nil, error_code.InvalidInput.Error(err)
// 	}

// 	json, err := c.GetJSONs()
// 	if err != nil {
// 		return nil, error_code.InvalidInput.Error(err)
// 	}
// 	json = json.MatchKeys(m.fields)

// 	data := json.MapString()
// 	if len(data) == 0 {
// 		return nil, error_code.InvalidInput.Text("没有可更新数据")
// 	}

// 	err = m.ActionUpdate(key, data)
// 	if err != nil {
// 		return nil, error_code.InvalidInput.Error(err)
// 	}

// 	return nil, nil
// }
