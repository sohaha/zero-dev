package model

import (
	"net/http"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zvalid"
	"github.com/zlsgo/zdb"
)

func (h *RestApi) Init(g *znet.Engine) {
	var (
		db *zdb.DB
		di = h.App.Di
	)

	zerror.Panic(di.Resolve(&db))

	_ = modelsBindRouter(g)
}

func Success(c *znet.Context, data interface{}, msg ...string) error {
	v := znet.ApiData{Data: data}
	if len(msg) > 0 {
		v.Msg = msg[0]
	}
	c.JSON(http.StatusOK, v)
	return nil
}

func ResultPages(rows ztype.Maps, pages zdb.Pages) ztype.Map {
	return ztype.Map{
		"items": rows,
		"page":  pages,
	}
}

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
