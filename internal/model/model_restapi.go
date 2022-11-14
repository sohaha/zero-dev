package model

import (
	"errors"
	"zlsapp/internal/error_code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
)

type RestApi struct {
	service.App
	db   *zdb.DB
	Path string
}

func NewRestApi() service.Router {
	return &RestApi{
		Path: "/model",
	}
}

func (h *RestApi) Init(g *znet.Engine) {
	zerror.Panic(h.App.Di.Resolve(&h.db))
	err := modelsBindRouter(g)
	if err != nil {
		zerror.Panic(zerror.With(err, "绑定模型接口路由失败"))
	}
}

func (m *Model) restApiInfo(key string, hasPrefix bool, fn ...func(b *builder.SelectBuilder) error) (ztype.Map, error) {
	return m.FindOne(func(b *builder.SelectBuilder) error {
		if key != "" && key != "0" {
			if hasPrefix {
				b.Where(b.EQ(m.Table.Name+"."+IDKey, key))
			} else {
				b.Where(b.EQ(IDKey, key))
			}
		}

		if len(fn) > 0 {
			return fn[0](b)
		}

		return nil
	}, false)
}

func (m *Model) restApiGetInfo(c *znet.Context) (interface{}, error) {
	key := c.GetParam("key")

	fields := GetViewFields(m, "info")
	finalFields, tmpFields, quote, with, withMany := getFinalFields(m, c, fields)

	info, err := m.restApiInfo(key, quote, func(b *builder.SelectBuilder) error {
		table := m.Table.Name
		for k, v := range with {
			m, ok := Get(v.Model)
			if !ok {
				return errors.New("关联模型(" + v.Model + ")不存在")
			}

			t := m.Table.Name
			asName := k
			b.JoinWithOption("", b.As(t, asName),
				asName+"."+v.Foreign+" = "+table+"."+v.Key,
			)
			if len(v.Fields) > 0 {
				finalFields = append(finalFields, zarray.Map(v.Fields, func(_ int, v string) string {
					return asName + "." + v
				})...)
			} else {
				finalFields = append(finalFields, asName+".*")
			}
		}
		b.Select(finalFields...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	for k, v := range withMany {
		m, ok := Get(v.Model)
		if !ok {
			return nil, errors.New("关联模型(" + v.Model + ")不存在")
		}
		key := info.Get(v.Key)
		if !key.Exists() {
			return nil, errors.New("字段(" + v.Key + ")不存在，无法关联模型(" + v.Model + ")")
		}

		row, _ := m.Find(func(b *builder.SelectBuilder) error {
			k := key.Value()
			switch val := k.(type) {
			case []interface{}:
				b.Where(b.In(v.Foreign, val...))
			default:
				b.Where(b.EQ(v.Foreign, val))
			}
			b.Select(m.GetFields()...)
			return nil
		}, false)
		_ = info.Set(k, row)
	}

	for _, v := range tmpFields {
		_ = info.Delete(v)
	}

	return info, nil

}

func (m *Model) restApiDelete(c *znet.Context) (interface{}, error) {
	key := c.GetParam("key")
	_, err := m.restApiInfo(key, false)
	if err != nil {
		return nil, err
	}

	if m.Options.SoftDeletes {
		_, err = m.DB.Update(m.Table.Name, map[string]interface{}{
			DeletedAtKey: ztime.Time().Unix(),
		}, func(b *builder.UpdateBuilder) error {
			b.Where(b.EQ(IDKey, key))
			return nil
		})
	} else {
		_, err = m.DB.Delete(m.Table.Name, func(b *builder.DeleteBuilder) error {
			b.Where(b.EQ(IDKey, key))
			return nil
		})
	}

	return nil, err
}

func (m *Model) restApiGetPage(c *znet.Context) (interface{}, error) {
	page, pagesize, err := GetPages(c)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	fields := GetViewFields(m, "lists")
	finalFields, tmpFields, quote, with, withMany := getFinalFields(m, c, fields)

	rows, pages, err := m.DB.Pages(m.Table.Name, page, pagesize, func(b *builder.SelectBuilder) error {
		deletedAtKey := DeletedAtKey
		idKey := IDKey
		if quote {
			table := m.Table.Name
			for k, v := range with {
				m, ok := Get(v.Model)
				if !ok {
					return errors.New("关联模型(" + v.Model + ")不存在")
				}

				t := m.Table.Name
				asName := k
				b.JoinWithOption("", b.As(t, asName),
					asName+"."+v.Foreign+" = "+table+"."+v.Key,
				)
				if len(v.Fields) > 0 {
					finalFields = append(finalFields, zarray.Map(v.Fields, func(_ int, v string) string {
						return asName + "." + v
					})...)
				} else {
					finalFields = append(finalFields, asName+".*")
				}
			}
			b.Desc(table + "." + IDKey)
			deletedAtKey = m.Table.Name + "." + DeletedAtKey
			idKey = m.Table.Name + "." + IDKey
		}

		b.Select(finalFields...)
		b.Desc(idKey)
		if m.Options.SoftDeletes {
			b.Where(b.EQ(deletedAtKey, 0))
		}
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

	return ResultPages(rows, pages), nil

}

func (m *Model) restApiCreate(c *znet.Context) (interface{}, error) {
	json, err := c.GetJSONs()
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	json = json.MatchKeys(m.fields)
	data := json.MapString()

	id, err := m.ActionCreate(data)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	return ztype.Map{"id": id}, nil
}

func (m *Model) restApiUpdate(c *znet.Context) (interface{}, error) {
	key := c.GetParam("key")
	_, err := m.restApiInfo(key, false)
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

	err = m.ActionUpdate(key, data)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	return nil, nil
}
