package model

import (
	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/schema"
)

// func parseRelation(m *Model, c *Column) {

// }

func parseColumn(m *Model, c *Column) {
	if c.ReadOnly {
		m.readOnlyKeys = append(m.readOnlyKeys, c.Name)
	}

	if c.Type == schema.JSON {
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

	parseValidRule(c)
	parseOptions(c)
}

// ParseJSON 解析模型
func ParseJSON(db *zdb.DB, json []byte) (m *Model, err error) {
	err = zjson.Unmarshal(json, &m)
	if err == nil {
		m.Raw = json
		m.DB = db
		m.readOnlyKeys = make([]string, 0)
		m.cryptKeys = make(map[string]cryptProcess, 0)
		m.afterProcess = make(map[string][]afterProcess, 0)
		m.beforeProcess = make(map[string][]beforeProcess, 0)

		// fillColumns(m)
		m.columnsKeys = zarray.Map(m.Columns, func(_ int, c *Column) string {
			parseColumn(m, c)
			// parseRelation(m, c)
			return c.Name
		})

		// m.relationKeys =
		// convertRelation(m)
	}
	return
}

func fillColumns(m *Model) {
	if m.Options.SoftDeletes {
		m.Columns = append(m.Columns, &Column{
			Name:     DeletedAtKey,
			Type:     schema.Int,
			Nullable: false,
			Comment:  "软删除时间",
		})
	}

	if m.Options.Timestamps {
		m.Columns = append(m.Columns, &Column{
			Name:    CreatedAtKey,
			Type:    schema.Time,
			Comment: "创建时间",
		}, &Column{
			Name:    UpdatedAtKey,
			Type:    schema.Time,
			Comment: "更新时间",
		})
	}
}
