package model

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/zlsgo/zdb"
)

// ParseJSON 解析模型
func ParseJSON(db *zdb.DB, json []byte) (m *Model, err error) {
	err = zjson.Unmarshal(json, &m)
	if err == nil {
		m.DB = db
		m.readOnlyKeys = make([]string, 0)
		m.cryptKeys = make(map[string]cryptProcess, 0)
		m.columnsKeys = zarray.Map(m.Columns, func(_ int, c *Column) string {
			if c.ReadOnly {
				m.readOnlyKeys = append(m.readOnlyKeys, c.Name)
			}
			if c.Crypt != "" {
				p, err := m.GetCryptProcess(c.Crypt)
				if err == nil {
					m.cryptKeys[c.Name] = p
				}
			}
			return c.Name
		})
	}
	return
}

func Add(db *zdb.DB, name string, json []byte) (m *Model, err error) {
	m, err = ParseJSON(db, json)
	if err == nil {
		name = strings.TrimSuffix(name, ".model.json")
		globalModels.Set(strings.Replace(name, "/", "-", -1), m)
	}
	return
}
