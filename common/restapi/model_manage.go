package restapi

import (
	"zlsapp/internal/error_code"
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type ManageRestApi struct {
	service.App
	Path string
}

func (h *ManageRestApi) Init(g *znet.Engine) {
	_ = g.GET("", h.lists)
	_ = g.GET("/:model", h.info)
	_ = g.GET("/:model/views", h.views)
}

// info 获取模型详情
func (h *ManageRestApi) info(c *znet.Context) error {
	modelName := c.GetParam("model")
	m, ok := parse.GetModel(modelName)
	if !ok {
		return error_code.InvalidInput.Text("模型不存在")
	}

	data := m.Raw
	data, _ = zjson.SetRawBytes([]byte(`{"code":0}`), "data", data)
	c.Byte(200, data)
	c.SetContentType(znet.ContentTypeJSON)
	return nil
}

// lists 获取模型列表
func (h *ManageRestApi) lists(c *znet.Context) (interface{}, error) {
	itmes := ztype.Maps{}
	parse.ModelsForEach(func(key string, m *parse.Modeler) bool {
		if m.Path == "" {
			return true
		}

		itmes = append(itmes, ztype.Map{
			"key":   key,
			"name":  m.Name,
			"table": m.Table.Name,
			"views": zarray.Keys(m.Views),
		})
		return true
	})

	return itmes, nil
}

// views 获取模型页面配置
func (h *ManageRestApi) views(c *znet.Context) (interface{}, error) {
	modelName := c.GetParam("model")
	m, ok := parse.GetModel(modelName)
	if !ok {
		return nil, error_code.InvalidInput.Text("模型不存在")
	}

	return m.GetView(), nil
}
