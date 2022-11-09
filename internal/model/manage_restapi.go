package model

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
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
}

// Get 获取模型列表
func (h *ManageRestApi) Get(c *znet.Context) (interface{}, error) {
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
