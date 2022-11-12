package model

import (
	"time"
	"zlsapp/internal/error_code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
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
			"views": zarray.Keys(m.Views),
		})
		return true
	})

	return itmes, nil
}

// info 获取模型详情
func (h *ManageRestApi) info(c *znet.Context) error {
	modelName := c.GetParam("model")
	key, _ := c.GetQuery("key")
	m, ok := globalModels.Get(modelName)
	if !ok {
		return error_code.InvalidInput.Text("模型不存在")
	}

	time.Sleep(time.Second * 1)
	data := m.Raw
	if key != "" {
		// json := zjson.ParseBytes(data).Get(key)
		// data = []byte(json.Raw())
		// title:= json.Get("title")
		// if !title.Exists() {
		// 	data = zjson.SetRawBytes(data, "data.", data)
		// }
		switch key {
		case "views":
			data, _ = zjson.SetBytes([]byte(`{}`), "", m.Views)
			zlog.Debug(m.Views)
		}
	}

	data, _ = zjson.SetRawBytes([]byte(`{"code":0}`), "data", data)
	c.Byte(200, data)
	c.SetContentType(znet.ContentTypeJSON)
	return nil
}
