package model

import (
	"strings"
	"zlsapp/internal/error_code"
	"zlsapp/internal/model/storage"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/znet"
)

var globalModels = zarray.NewHashMap[string, *Model]()

func Get(name string) (*Model, bool) {
	return globalModels.Get(name)
}

func Add(name string, json []byte, bindStorage func(*Model) (storage.Storageer, error), force ...bool) (m *Model, err error) {
	err = ValidateModelSchema(json)
	if err != nil {
		err = zerror.With(err, "模型("+name+")验证失败")
		return
	}
	m, err = ParseJSON(json)
	if err == nil {
		m.Storage, err = bindStorage(m)
		if err != nil {
			return
		}
		name = strings.TrimSuffix(name, ".model.json")
		name = strings.Replace(name, "/", "-", -1)
		if _, ok := globalModels.Get(name); ok && !(len(force) > 0 && force[0]) {
			return nil, ErrModuleAlreadyExists
		}
		globalModels.Set(name, m)
	}
	return
}

func modelsBindRouter(g *znet.Engine) error {
	g.Use(func(c *znet.Context) error {
		name := c.GetParam("model")
		m, ok := globalModels.Get(name)
		if !ok {
			return error_code.NotFound.Text("模型不存在")
		}

		c.WithValue("model", m)
		c.Next()
		return nil
	})

	g.GET(":model", func(c *znet.Context) (interface{}, error) {
		return c.MustValue("model").(*Model).restApiGetPage(c)
	})

	g.GET(":model/:key", func(c *znet.Context) (interface{}, error) {
		return c.MustValue("model").(*Model).restApiGetInfo(c)
	})

	g.POST(":model", func(c *znet.Context) (interface{}, error) {
		return c.MustValue("model").(*Model).restApiCreate(c)
	})

	g.PATCH(":model/:key", func(c *znet.Context) (interface{}, error) {
		return c.MustValue("model").(*Model).restApiUpdate(c)
	})

	g.DELETE(":model/:key", func(c *znet.Context) (interface{}, error) {
		return c.MustValue("model").(*Model).restApiDelete(c)
	})

	return nil
}
