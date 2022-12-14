package parse

import (
	"errors"
	"strings"
	"zlsapp/common"
	"zlsapp/common/hashid"
	"zlsapp/internal/error_code"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zvalid"
	"github.com/zlsgo/zdb"
)

type Api struct {
	Method string
	Path   string
	Handle znet.Handler
	Public bool
}

type ApiKeyType string

const (
	ApiKeyPages  ApiKeyType = "pages"
	ApiKeyQuery  ApiKeyType = "query"
	ApiKeyCreate ApiKeyType = "create"
	ApiKeyUpdate ApiKeyType = "update"
	ApiKeyDelete ApiKeyType = "delete"
)

func JudgeRouters(m *Modeler, t ApiKeyType) bool {
	if _, ok := m.apis[t]; !ok {
		return false
	}

	return true
}

func resolverApi(m *Modeler) {
	m.apis = make(map[ApiKeyType]Api, 0)
	apiOptions := ztype.New(m.Options.Api)
	if !apiOptions.Bool() {
		return
	}
	opts := apiOptions.Map()
	if _, isBool := apiOptions.Value().(bool); isBool {
		opts = ztype.Map{
			"query": ztype.Map{
				"public": true,
			},
			"pages": ztype.Map{
				"public": true,
			},
		}
	}

	zlog.Debug(m.Name, opts)
	for k := range opts {
		v := opts.Get(k)
		api := Api{
			Path:   "/api/" + m.Alias,
			Public: v.Get("public").Bool(),
		}
		// switch k {
		// case "query":
		// 	api.Path = api.Path + "/:key"
		// 	api.Method = "GET"
		// 	api.Handle = func(c *znet.Context) (interface{}, error) {
		// 		return RestapiGetInfo(c, m, []string{}, []string{})
		// 	}
		// case "pages":
		// 	api.Method = "GET"
		// 	api.Handle = func(c *znet.Context) (interface{}, error) {
		// 		return RestapiGetPage(c, m, ztype.Map{}, []string{}, []string{})
		// 	}
		// case "create":
		// 	api.Method = "POST"
		// case "update":
		// 	api.Method = "PUT"
		// 	api.Path = api.Path + "/:key"
		// case "delete":
		// 	api.Path = api.Path + "/:key"
		// 	api.Method = "DELETE"
		// default:
		// 	zlog.Error("api method not found", k)
		// 	continue
		// }
		// if api.Handle == nil {
		// 	continue
		// }
		m.apis[ApiKeyType(k)] = api
	}

}

type PageData struct {
	Page  PageInfo   `json:"page"`
	Items ztype.Maps `json:"items"`
}

func GetPages(c *znet.Context) (page, pagesize int, err error) {
	rule := c.ValidRule().IsNumber().MinInt(1)
	err = zvalid.Batch(
		zvalid.BatchVar(&page, c.Valid(rule, "page", "??????").Customize(func(rawValue string, err error) (string, error) {
			if err != nil || rawValue == "" {
				return "1", err
			}
			return rawValue, nil
		})),

		zvalid.BatchVar(&pagesize, c.Valid(rule, "pagesize", "??????").MaxInt(1000).Customize(func(rawValue string, err error) (string, error) {
			if err != nil || rawValue == "" {
				return "10", err
			}
			return rawValue, nil
		})),
	)
	return
}

func restApiInfo(m *Modeler, key string, filter ztype.Map, fn ...StorageOptionFn) (ztype.Map, error) {
	if key != "" && key != "0" {
		filter[IDKey] = key
	}

	row, err := FindOne(m, filter, fn...)
	if (err != nil && err == zdb.ErrNotFound) || row.IsEmpty() {
		err = errors.New("???????????????")
	}
	return row, err
}

func getRestapiKey(c *znet.Context, m *Modeler) (string, error) {
	key := c.GetParam("key")

	if m.Options.CryptID {
		id, err := hashid.DecryptID(m.hashid, key)
		if err != nil {
			return "", errors.New("ID ????????????")
		}
		key = ztype.ToString(id)
	}
	return key, nil
}

func RestapiGetInfo(c *znet.Context, m *Modeler, filter ztype.Map, fields []string, withFilds []string) (interface{}, error) {
	key, err := getRestapiKey(c, m)
	if err != nil {
		return nil, err
	}
	finalFields, tmpFields, with, withMany := getFinalFields(m, c, fields, withFilds)

	info, err := restApiInfo(m, key, filter, func(so *StorageOptions) error {
		table := m.Table.Name
		for k, v := range with {
			m, ok := GetModel(v.Model)
			if !ok {
				return errors.New("????????????(" + v.Model + ")?????????")
			}

			t := m.Table.Name
			asName := k
			so.Join = append(so.Join, StorageJoin{
				Table:  t,
				As:     asName,
				Option: v.Join,
				Expr:   asName + "." + v.Foreign + " = " + table + "." + v.Key,
			})

			if len(v.Fields) > 0 {
				finalFields = append(finalFields, zarray.Map(v.Fields, func(_ int, v string) string {
					return asName + "." + v
				})...)
			} else {
				finalFields = append(finalFields, asName+".*")
			}
		}
		so.Fields = finalFields
		return nil
	})

	if err != nil {
		return nil, err
	}

	for k, v := range withMany {
		m, ok := GetModel(v.Model)
		if !ok {
			return nil, errors.New("????????????(" + v.Model + ")?????????")
		}
		key := info.Get(v.Key)
		if !key.Exists() {
			return nil, errors.New("??????(" + v.Key + ")??????????????????????????????(" + v.Model + ")")
		}

		rows, _ := Find(m, ztype.Map{
			v.Foreign: key.Value(),
		}, func(so *StorageOptions) error {
			if len(v.Fields) > 0 {
				so.Fields = v.Fields
			}
			return nil
		})

		_ = info.Set(k, rows)
	}

	for _, v := range tmpFields {
		s := strings.SplitN(v, ".", 2)
		if len(s) == 2 {
			_ = info.Delete(s[1])
		} else {
			_ = info.Delete(v)
		}
	}

	if m.Options.CryptID && zarray.Contains(finalFields, IDKey) {
		id := info.Get(IDKey)
		hid, err := hashid.EncryptID(m.hashid, id.Int64())
		if err != nil {
			return nil, zerror.With(err, "?????? ID ??????")
		}
		_ = info.Set(IDKey, hid)
	}

	return info, nil
}

func RestapiGetPage(c *znet.Context, m *Modeler, filter ztype.Map, fields []string, withFilds []string, fn ...func(so *StorageOptions) error) (*PageData, error) {
	page, pagesize, err := GetPages(c)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	finalFields, tmpFields, with, withMany := getFinalFields(m, c, fields, withFilds)

	rows, pageInfo, err := Pages(m, page, pagesize, filter, func(so *StorageOptions) error {
		so.OrderBy = map[string]int8{m.Table.Name + "." + IDKey: -1}
		if len(with) > 0 {
			table := m.Table.Name
			for k, v := range with {
				m, ok := GetModel(v.Model)
				if !ok {
					return errors.New("????????????(" + v.Model + ")?????????")
				}

				t := m.Table.Name
				asName := k
				so.Join = append(so.Join, StorageJoin{
					Table:  t,
					As:     k,
					Option: v.Join,
					Expr:   asName + "." + v.Foreign + " = " + table + "." + v.Key,
				})

				if len(v.Fields) > 0 {
					finalFields = append(finalFields, zarray.Map(v.Fields, func(_ int, v string) string {
						// TODO ???????????????,???????????????
						return asName + "." + v
					})...)
				} else {
					finalFields = append(finalFields, asName+".*")
				}
			}
		}

		so.Fields = finalFields
		for _, f := range fn {
			if err = f(so); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// TODO dev
	_ = withMany
	for _, info := range rows {
		if m.Options.CryptID && zarray.Contains(finalFields, IDKey) {
			id := info.Get(IDKey)
			hid, err := hashid.EncryptID(m.hashid, id.Int64())
			if err != nil {
				return nil, zerror.With(err, "?????? ID ??????")
			}
			_ = info.Set(IDKey, hid)
		}
		for _, v := range tmpFields {
			_ = info.Delete(v)
		}
	}

	return &PageData{
		Items: rows,
		Page:  pageInfo,
	}, nil

}

func RestapiCreate(c *znet.Context, m *Modeler) (interface{}, error) {
	json, err := c.GetJSONs()
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	json = json.MatchKeys(m.fields)
	data := json.MapString()

	uid := common.GetUID(c)
	id, err := Insert(m, data, uid)

	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	return ztype.Map{"id": id}, nil
}

func RestapiDelete(c *znet.Context, m *Modeler, filter ztype.Map) (interface{}, error) {
	key, err := getRestapiKey(c, m)
	if err != nil {
		return nil, err
	}

	_, err = restApiInfo(m, key, filter)
	if err != nil {
		return nil, err
	}

	_, err = Delete(m, key)

	return nil, err
}

func RestapiUpdate(c *znet.Context, m *Modeler, filter ztype.Map) (interface{}, error) {
	key, err := getRestapiKey(c, m)
	if err != nil {
		return nil, err
	}

	_, err = restApiInfo(m, key, filter)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	json, err := c.GetJSONs()
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}
	json = json.MatchKeys(m.fields)

	data := json.MapString()
	if len(data) == 0 {
		return nil, error_code.InvalidInput.Text("?????????????????????")
	}

	_, err = Update(m, key, data)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	return nil, nil
}
