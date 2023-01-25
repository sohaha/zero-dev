package api

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"zlsapp/common"
	"zlsapp/common/hashid"
	"zlsapp/internal/error_code"
	"zlsapp/internal/parse"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
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

// type ApiKeyType string
//
// const (
// 	ApiKeyPages  ApiKeyType = "pages"
// 	ApiKeyQuery  ApiKeyType = "query"
// 	ApiKeyCreate ApiKeyType = "create"
// 	ApiKeyUpdate ApiKeyType = "update"
// 	ApiKeyDelete ApiKeyType = "delete"
// )
//
// func JudgeRouters(m *parse.Modeler, t ApiKeyType) bool {
// 	if _, ok := m.apis[t]; !ok {
// 		return false
// 	}
//
// 	return true
// }

// func resolverApi(m *parse.Modeler) {
// 	m.apis = make(map[ApiKeyType]Api, 0)
// 	apiOptions := ztype.New(m.Options.Api)
// 	if !apiOptions.Bool() {
// 		return
// 	}
// 	opts := apiOptions.Map()
// 	if _, isBool := apiOptions.Value().(bool); isBool {
// 		opts = ztype.Map{
// 			"query": ztype.Map{
// 				"public": true,
// 			},
// 			"pages": ztype.Map{
// 				"public": true,
// 			},
// 		}
// 	}
//
// 	zlog.Debug(m.Name, opts)
// 	for k := range opts {
// 		v := opts.Get(k)
// 		api := Api{
// 			Path:   "/api/" + m.Alias,
// 			Public: v.Get("public").Bool(),
// 		}
// 		// switch k {
// 		// case "query":
// 		// 	api.Path = api.Path + "/:key"
// 		// 	api.Method = "GET"
// 		// 	api.Handle = func(c *znet.Context) (interface{}, error) {
// 		// 		return RestapiGetInfo(c, m, []string{}, []string{})
// 		// 	}
// 		// case "pages":
// 		// 	api.Method = "GET"
// 		// 	api.Handle = func(c *znet.Context) (interface{}, error) {
// 		// 		return RestapiGetPage(c, m, ztype.Map{}, []string{}, []string{})
// 		// 	}
// 		// case "create":
// 		// 	api.Method = "POST"
// 		// case "update":
// 		// 	api.Method = "PUT"
// 		// 	api.Path = api.Path + "/:key"
// 		// case "delete":
// 		// 	api.Path = api.Path + "/:key"
// 		// 	api.Method = "DELETE"
// 		// default:
// 		// 	zlog.Error("api method not found", k)
// 		// 	continue
// 		// }
// 		// if api.Handle == nil {
// 		// 	continue
// 		// }
// 		m.apis[ApiKeyType(k)] = api
// 	}
//
// }

type PageData struct {
	Page  parse.PageInfo `json:"page"`
	Items ztype.Maps     `json:"items"`
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

func restApiInfo(m *parse.Modeler, key string, filter ztype.Map, fn ...parse.StorageOptionFn) (ztype.Map, error) {
	if key != "" && key != "0" {
		filter[parse.IDKey] = key
	}

	row, err := parse.FindOne(m, filter, fn...)
	if (err != nil && err == zdb.ErrNotFound) || row.IsEmpty() {
		err = errors.New("记录不存在")
	}
	return row, err
}

func getRestapiKey(c *znet.Context, m *parse.Modeler) (string, error) {
	key := c.GetParam("key")

	if m.Options.CryptID {
		id, err := hashid.DecryptID(m.Hashid, key)
		if err != nil {
			return "", errors.New("ID 解密失败")
		}
		key = ztype.ToString(id)
	}
	return key, nil
}

func RestapiGetInfo(c *znet.Context, m *parse.Modeler, filter ztype.Map, fields []string, withFilds []string) (interface{}, error) {
	key, err := getRestapiKey(c, m)
	if err != nil {
		return nil, err
	}
	finalFields, tmpFields, with, withMany := getFinalFields(m, c, fields, withFilds)

	info, err := restApiInfo(m, key, filter, func(so *parse.StorageOptions) error {
		table := m.Table.Name
		for k, v := range with {
			m, ok := parse.GetModel(v.Model)
			if !ok {
				return errors.New("关联模型(" + v.Model + ")不存在")
			}

			t := m.Table.Name
			asName := k
			so.Join = append(so.Join, parse.StorageJoin{
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
		m, ok := parse.GetModel(v.Model)
		if !ok {
			return nil, errors.New("关联模型(" + v.Model + ")不存在")
		}
		key := info.Get(v.Key)
		if !key.Exists() {
			return nil, errors.New("字段(" + v.Key + ")不存在，无法关联模型(" + v.Model + ")")
		}

		rows, _ := parse.Find(m, ztype.Map{
			v.Foreign: key.Value(),
		}, func(so *parse.StorageOptions) error {
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

	if m.Options.CryptID && zarray.Contains(finalFields, parse.IDKey) {
		id := info.Get(parse.IDKey)
		hid, err := hashid.EncryptID(m.Hashid, id.Int64())
		if err != nil {
			return nil, zerror.With(err, "加密 ID 失败")
		}
		_ = info.Set(parse.IDKey, hid)
	}

	return info, nil
}

func RestapiGetPage(c *znet.Context, m *parse.Modeler, filter ztype.Map, fields []string, withFilds []string, fn ...func(so *parse.StorageOptions) error) (*PageData, error) {
	page, pagesize, err := GetPages(c)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	finalFields, tmpFields, with, withMany := getFinalFields(m, c, fields, withFilds)

	rows, pageInfo, err := parse.Pages(m, page, pagesize, filter, func(so *parse.StorageOptions) error {
		so.OrderBy = map[string]int8{m.Table.Name + "." + parse.IDKey: -1}
		if len(with) > 0 {
			table := m.Table.Name
			for k, v := range with {
				m, ok := parse.GetModel(v.Model)
				if !ok {
					return errors.New("关联模型(" + v.Model + ")不存在")
				}

				t := m.Table.Name
				asName := k
				so.Join = append(so.Join, parse.StorageJoin{
					Table:  t,
					As:     k,
					Option: v.Join,
					Expr:   asName + "." + v.Foreign + " = " + table + "." + v.Key,
				})

				if len(v.Fields) > 0 {
					finalFields = append(finalFields, zarray.Map(v.Fields, func(_ int, v string) string {
						// TODO 修改了字段,这是暂时的
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
		if m.Options.CryptID && zarray.Contains(finalFields, parse.IDKey) {
			id := info.Get(parse.IDKey)
			hid, err := hashid.EncryptID(m.Hashid, id.Int64())
			if err != nil {
				return nil, zerror.With(err, "加密 ID 失败")
			}
			_ = info.Set(parse.IDKey, hid)
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

func RestapiCreate(c *znet.Context, m *parse.Modeler) (interface{}, error) {
	json, err := c.GetJSONs()
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	json = json.MatchKeys(m.Fields)
	data := json.MapString()

	uid := common.GetUID(c)
	id, err := parse.Insert(m, data, uid)

	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	return ztype.Map{"id": id}, nil
}

func RestapiDelete(c *znet.Context, m *parse.Modeler, filter ztype.Map) (interface{}, error) {
	key, err := getRestapiKey(c, m)
	if err != nil {
		return nil, err
	}

	_, err = restApiInfo(m, key, filter)
	if err != nil {
		return nil, err
	}

	_, err = parse.Delete(m, key)

	return nil, err
}

func RestapiUpdate(c *znet.Context, m *parse.Modeler, filter ztype.Map) (interface{}, error) {
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
	json = json.MatchKeys(m.Fields)

	data := json.MapString()
	if len(data) == 0 {
		return nil, error_code.InvalidInput.Text("没有可更新数据")
	}

	_, err = parse.Update(m, key, data)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	return nil, nil
}

type UploadOption struct {
	Key      string
	Dir      string
	MimeType []string
	MaxSize  int64
}

func RestapiUpload(c *znet.Context, m *parse.Modeler, opt ...func(o *UploadOption)) (interface{}, error) {
	o := UploadOption{
		Key:     "file",
		MaxSize: 1024 * 1024 * 2,
	}

	files, err := c.FormFiles(o.Key)
	if err != nil {
		return nil, zerror.InvalidInput.Wrap(err, "上传失败")
	}

	uploadDir := zfile.RealPathMkdir("resource/upload/"+m.Alias+"/"+o.Dir, true)

	// dir := c.DefaultFormOrQuery("dir", "")
	// dir = zfile.RealPathMkdir(uploadDir+dir, true)
	// if !zfile.IsSubPath(dir, uploadDir) {
	// 	return nil, error_code.InvalidInput.Text("非法存储目录")
	// }

	uploads := make(map[string]*multipart.FileHeader, len(files))
	buf := bytes.NewBuffer(nil)

	for _, v := range files {
		f, err := v.Open()
		if err != nil {
			return nil, zerror.InvalidInput.Wrap(err, "文件读取失败")
		}

		if _, err := io.Copy(buf, f); err != nil {
			if err != nil {
				return nil, zerror.InvalidInput.Wrap(err, "文件读取失败")
			}
		}

		f.Close()

		b := buf.Bytes()
		if len(b) > int(o.MaxSize) {
			return nil, zerror.InvalidInput.Wrap(err, "文件大小超出限制")
		}

		mt := zfile.GetMimeType(v.Filename, b)
		n := strings.Split(mt, "/")
		if len(n) < 2 {
			return nil, zerror.InvalidInput.Wrap(err, "文件类型错误")
		}

		if len(o.MimeType) > 0 {
			ok := false
			for _, v := range o.MimeType {
				if v == mt || v == n[1] {
					ok = true
					break
				}
			}

			if !ok {
				return nil, error_code.InvalidInput.Text("不支持的文件类型")
			}
		}

		ext := filepath.Ext(v.Filename)
		if ext == "" {
			if len(n) > 1 {
				ext = "." + n[len(n)]
			}
		}

		id := zstring.Md5Byte(b) + ext
		uploads[id] = v

		buf.Reset()
	}

	for n, f := range uploads {
		err = c.SaveUploadedFile(f, uploadDir+n)
		if err != nil {
			return nil, error_code.ServerError.Text("文件保存失败", err)
		}
	}

	return zarray.Map(zarray.Keys(uploads), func(_ int, p string) string {
		return "/" + zfile.SafePath(uploadDir+p, "resource")
	}), nil

}
