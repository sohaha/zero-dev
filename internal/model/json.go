package model

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/schema"
)

// ParseJSON 解析模型
func ParseJSON(db *zdb.DB, json []byte) (m *Model, err error) {
	err = zjson.Unmarshal(json, &m)
	if err == nil {
		m.DB = db
		m.readOnlyKeys = make([]string, 0)
		m.cryptKeys = make(map[string]cryptProcess, 0)
		m.afterProcess = make(map[string][]afterProcess, 0)
		m.beforeProcess = make(map[string][]beforeProcess, 0)
		m.columnsKeys = zarray.Map(m.Columns, func(_ int, c *Column) string {
			if c.ReadOnly {
				m.readOnlyKeys = append(m.readOnlyKeys, c.Name)
			}
			if c.Type == string(schema.JSON) {
				if len(c.Before) == 0 {
					c.Before = []string{"json"}
				}
				if len(c.After) == 0 {
					c.After = []string{"json"}
				}
			}
			if c.Crypt != "" {
				p, err := m.GetCryptProcess(c.Crypt)
				if err == nil {
					m.cryptKeys[c.Name] = p
				}
			}

			if len(c.Before) > 0 {
				ps, err := m.GetBeforeProcess(c.Before)
				if err == nil {
					m.beforeProcess[c.Name] = ps
				}
			}

			if len(c.After) > 0 {
				ps, err := m.GetAfterProcess(c.Before)
				if err == nil {
					m.afterProcess[c.Name] = ps
				}
			}
			return c.Name
		})
	}
	return
}

func Add(db *zdb.DB, name string, json []byte, force bool) (m *Model, err error) {
	m, err = ParseJSON(db, json)
	if err == nil {
		name = strings.TrimSuffix(name, ".model.json")
		name = strings.Replace(name, "/", "-", -1)
		if _, ok := globalModels.Get(name); ok && !force {
			return nil, errors.New("model already exists")
		}
		globalModels.Set(name, m)
	}
	return
}
