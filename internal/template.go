package app

import (
	"net/http"
	"strings"
	"zlsapp/internal/account"
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/jet"
)

func InitTemplate() *service.Template {
	return &service.Template{
		DIR: "app/templates",
	}
}

func bindModelTemplate(r *znet.Engine, di zdi.Invoker) {
	dir := "app/templates"
	if !zfile.DirExist(dir) {
		return
	}

	var conf *service.Conf
	_ = di.Resolve(&conf)

	j := jet.New(r, dir, func(o *jet.Options) {})

	_ = j.Load()

	injectionTemplate(j)
	r.SetTemplate(j)

	var mapping ztype.Map

	if zfile.FileExist(dir + "/mapping.json") {
		f, _ := zfile.ReadFile(dir + "/mapping.json")
		_ = zjson.Unmarshal(f, &mapping)
	}

	_, hasHome := zarray.Find(zarray.Keys(mapping), func(_ int, k string) bool { return k == "/" })
	if !hasHome && j.Exists("index") {
		r.GET("/", func(c *znet.Context) {
			c.Template(http.StatusOK, "index", inlayTemplateArgs(c, conf, nil))
		})
	}

	for k := range mapping {
		v := mapping.Get(k)
		k := strings.TrimLeft(k, "/")
		r.GET("/"+k, func(c *znet.Context) {
			template := v.String()
			if !j.Exists(template) {
				c.String(400, "模板不存在")
				return
			}

			c.Template(http.StatusOK, template, inlayTemplateArgs(c, conf, nil))
		}, znet.WrapFirstMiddleware(func(c *znet.Context) {
			c.WithValue(account.DisabledAuthKey, true)
			c.Next()
		}))
	}
}

func inlayTemplateArgs(c *znet.Context, conf *service.Conf, m ztype.Map) ztype.Map {
	if m == nil {
		m = ztype.Map{}
	}
	m["conf"] = conf.Core().Get("app")
	m["path"] = c.Request.URL.Path
	m["host"] = c.Host()
	m["params"] = c.GetAllParam()
	m["query"] = func(key string) string {
		return c.DefaultQuery("key", "")
	}
	return m
}

func injectionTemplate(j *jet.Engine) {
	j.AddFunc("FindOne", func(model string, id interface{}, data ...ztype.Map) ztype.Map {
		m, ok := parse.GetModel(model)
		if !ok {
			return ztype.Map{}
		}
		filter := ztype.Map{}
		if len(data) > 0 {
			filter = data[0]
		}
		i := ztype.ToString(id)
		if i != "" && i != "0" {
			filter[parse.IDKey] = id
		}
		item, _ := parse.FindOne(m, filter, func(so *parse.StorageOptions) error {
			so.Fields = m.GetFields()
			return nil
		})
		return item
	})

	j.AddFunc("Find", func(model string, data ...ztype.Map) ztype.Maps {
		m, ok := parse.GetModel(model)
		if !ok {
			return ztype.Maps{}
		}

		filter := ztype.Map{}
		if len(data) > 0 {
			filter = data[0]
		}

		items, _ := parse.Find(m, filter, func(so *parse.StorageOptions) error {
			so.Fields = m.GetFields()
			return nil
		})

		return items
	})

	j.AddFunc("QueryGroup", func(model string, field string, data ...ztype.Map) []string {
		m, ok := parse.GetModel(model)
		if !ok {
			return []string{}
		}

		filter := ztype.Map{}
		if len(data) > 0 {
			filter = data[0]
		}

		items, _ := parse.Find(m, filter, func(so *parse.StorageOptions) error {
			so.Fields = []string{field}
			so.GroupBy = []string{field}
			return nil
		})

		return zarray.Map(items, func(_ int, m ztype.Map) string {
			return m.Get(field).String()
		})
	})

	j.AddFunc("Pages", func(model string, page, pagesize interface{}, filter ztype.Map) ztype.Map {
		m, ok := parse.GetModel(model)
		if !ok {
			return ztype.Map{}
		}
		items, p, _ := parse.Pages(m, ztype.ToInt(page), ztype.ToInt(pagesize), filter, func(so *parse.StorageOptions) error {
			so.Fields = m.GetFields("content")
			return nil
		})

		return ztype.Map{
			"items": items,
			"page":  p,
		}
	})
	{

		j.AddFunc("FAQs", func(category string, data ...ztype.Map) ztype.Maps {
			m, ok := parse.GetModel("website-faq")
			if !ok {
				return ztype.Maps{}
			}

			filter := ztype.Map{}
			if len(data) > 0 {
				filter = data[0]
			}
			if category != "" {
				filter["category"] = category
			}

			items, _ := parse.Find(m, filter, func(so *parse.StorageOptions) error {
				so.Fields = []string{parse.IDKey, "title"}
				so.Limit = 6
				return nil
			})

			return items
		})

		j.AddFunc("FAQ", func(id string, data ...ztype.Map) ztype.Map {
			m, ok := parse.GetModel("website-faq")
			if !ok {
				return ztype.Map{}
			}

			filter := ztype.Map{
				parse.IDKey: id,
			}
			item, _ := parse.FindOne(m, filter, func(so *parse.StorageOptions) error {
				return nil
			})

			return item
		})

		j.AddFunc("News", func(page, pagesize interface{}) ztype.Map {
			m, ok := parse.GetModel("website-news")
			if !ok {
				return ztype.Map{}
			}
			filter := ztype.Map{}
			items, p, _ := parse.Pages(m, ztype.ToInt(page), ztype.ToInt(pagesize), filter, func(so *parse.StorageOptions) error {
				so.Fields = m.GetFields("content")
				return nil
			})

			return ztype.Map{
				"items": items,
				"page":  p,
			}
		})
	}
}
