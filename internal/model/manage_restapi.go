package model

import (
	"zlsapp/internal/error_code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
)

type ManageRestApi struct {
	service.App
	db   *zdb.DB
	Path string
}

func NewManageRestApi() service.Router {
	return &ManageRestApi{
		Path: "/manage/model",
	}
}

func (h *ManageRestApi) Init(g *znet.Engine) {
	zerror.Panic(h.App.Di.Resolve(&h.db))
	_ = g.GET("", h.lists)
	_ = g.GET("/:model", h.info)
}

// lists 获取模型列表
func (h *ManageRestApi) lists(c *znet.Context) (interface{}, error) {
	itmes := ztype.Maps{}
	globalModels.ForEach(func(key string, m *Model) bool {
		if m.Path == "" {
			return true
		}

		itmes = append(itmes, ztype.Map{
			"key":   key,
			"name":  m.Name,
			"table": m.Table.Name,
		})
		return true
	})

	return itmes, nil
}

// info 获取模型详情
func (h *ManageRestApi) info(c *znet.Context) error {
	modelName := c.GetParam("model")
	m, ok := globalModels.Get(modelName)
	if !ok {
		return error_code.InvalidInput.Text("模型不存在")
	}

	data, _ := zjson.SetRawBytes([]byte(`{"code":0}`), "data", m.Raw)
	data = zjson.Ugly(data)

	c.Byte(200, data)
	c.SetContentType(znet.ContentTypeJSON)
	return nil
}
