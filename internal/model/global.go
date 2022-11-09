package model

import (
	"errors"
	"strings"
	"zlsapp/internal/error_code"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
	"github.com/zlsgo/zdb"
)

var globalModels = zarray.NewHashMap[string, *Model]()

func Get(name string) (*Model, bool) {
	return globalModels.Get(name)
}

func Add(db *zdb.DB, name string, json []byte, force bool) (m *Model, err error) {
	m, err = ParseJSON(db, json)
	if err == nil {
		name = strings.TrimSuffix(name, ".model.json")
		name = strings.Replace(name, "/", "-", -1)
		if _, ok := globalModels.Get(name); ok && !force {
			return nil, errors.New("model already exists")
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
