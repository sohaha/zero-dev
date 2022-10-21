package model

import (
	"net/http"
	"zlsapp/app/error_code"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zvalid"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/builder"
)

type RestApi struct {
	service.App
	Path string
}

func NewRestApi() service.Router {
	return &RestApi{
		Path: "/api",
	}
}

func (h *RestApi) Init(g *znet.Engine) {
	di := h.App.Di

	var db *zdb.DB

	zerror.Panic(di.Resolve(&db))

	// g := r.Group("/api")
	// allModel := make([]*Model,0,globalModels.Len())
	globalModels.ForEach(func(s string, m *Model) bool {
		zlog.Debug(s, m.Name)
		g.Handle("GET", s, func(c *znet.Context) error {
			page, pagesize, err := GetPages(c)
			if err != nil {
				return error_code.InvalidInput.Error(err)
			}
			rows, pages, _ := db.Pages(m.Table.Name, page, pagesize, func(b *builder.SelectBuilder) error {
				b.Desc(IDKey)

				fields := m.columnsKeys
				fields = append(fields, IDKey)

				if m.Options.Timestamps {
					fields = append(fields, CreatedAtKey, UpdatedAtKey)
				}
				b.Select(fields...)

				zlog.Debug(fields)
				return nil
			})

			return Success(c, ResultPages(rows, pages))
		})
		return true
	})
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
