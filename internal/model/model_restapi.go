package model

import (
	"zlsapp/internal/error_code"
	"zlsapp/service"

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

func (m *Model) restApiInfo(key string, fn ...func(b *builder.SelectBuilder) error) (ztype.Map, error) {
	return m.FindOne(func(b *builder.SelectBuilder) error {
		if key != "" {
			b.Where(b.EQ(IDKey, key))
		}

		if len(fn) > 0 {
			return fn[0](b)
		}

		return nil
	}, false)
}

func (m *Model) restApiGetInfo(c *znet.Context) error {
	key := c.GetParam("key")
	fields := GetRequestFields(c, m)

	row, err := m.restApiInfo(key, func(b *builder.SelectBuilder) error {
		b.Select(fields...)
		return nil
	})
	if err != nil {
		return error_code.InvalidInput.Error(err)
	}

	return Success(c, row)
}

func (m *Model) restApiDelete(c *znet.Context) error {
	key := c.GetParam("key")
	_, err := m.restApiInfo(key)
	if err != nil {
		return error_code.InvalidInput.Error(err)
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

	if err != nil {
		return err
	}

	return Success(c, nil)

}

func (m *Model) restApiGetPage(c *znet.Context) error {
	page, pagesize, err := GetPages(c)
	if err != nil {
		return error_code.InvalidInput.Error(err)
	}

	fields := GetRequestFields(c, m)

	rows, pages, err := m.DB.Pages(m.Table.Name, page, pagesize, func(b *builder.SelectBuilder) error {
		b.Desc(IDKey)
		b.Select(fields...)
		if m.Options.SoftDeletes {
			b.Where(b.EQ(DeletedAtKey, 0))
		}
		return nil
	})

	if err != nil {
		return err
	}

	return Success(c, ResultPages(rows, pages))

}

func (m *Model) restApiCreate(c *znet.Context) error {
	json, err := c.GetJSONs()
	if err != nil {
		return error_code.InvalidInput.Error(err)
	}

	json = json.MatchKeys(m.columnsKeys)

	data := json.MapString()

	id, err := m.ActionCreate(data)
	if err != nil {
		return error_code.InvalidInput.Error(err)
	}

	return Success(c, ztype.Map{"id": id})
}

func (m *Model) restApiUpdate(c *znet.Context) error {
	key := c.GetParam("key")
	_, err := m.restApiInfo(key)
	if err != nil {
		return error_code.InvalidInput.Error(err)
	}

	json, err := c.GetJSONs()
	if err != nil {
		return error_code.InvalidInput.Error(err)
	}
	json = json.MatchKeys(m.columnsKeys)

	data := json.MapString()
	err = m.ActionUpdate(key, data)
	if err != nil {
		return error_code.InvalidInput.Error(err)
	}

	return Success(c, data)
}
