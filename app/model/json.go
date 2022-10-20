package model

import (
	"github.com/sohaha/zlsgo/zjson"
	"github.com/zlsgo/zdb"
)

// 解析
func ParseJSON(db *zdb.DB, json []byte) (m *Model, err error) {
	err = zjson.Unmarshal(json, &m)
	if err == nil {
		m.DB = db
	}
	return
	// j := zjson.ParseBytes(json)

	// columns := make([]Column, 0)
	// j.Get("columns").ForEach(func(key, value zjson.Res) bool {
	// 	zlog.Debug(key, value)
	// 	column := Column{
	// 		Name:        value.Get("name").String(),
	// 		Comment:     value.Get("comment").String(),
	// 		Type:        value.Get("type").String(),
	// 		Size:        value.Get("size").Uint(),
	// 		Nullable:    value.Get("nullable").Bool(),
	// 		Label:       value.Get("label").String(),
	// 		Enum:        nil,
	// 		Unique:      nil,
	// 		Index:       nil,
	// 		Validations: []validations{},
	// 	}

	// 	defaultValue := value.Get("default")
	// 	if defaultValue.Exists() {
	// 		column.Default = defaultValue.Value()
	// 	}

	// 	validations := value.Get("validations")
	// 	if validations.Exists() {
	// 		_ = zjson.Unmarshal(validations.String(), &column.Validations)
	// 	}
	// 	columns = append(columns, column)
	// 	return true
	// })
	// _ = columns
}
