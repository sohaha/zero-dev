package parse

import (
	"zlsapp/internal/parse/jsonschema"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/zlsgo/zdb/schema"
)

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
}

// ParseModel 解析模型
func ParseModel(json []byte) (m *Model, err error) {
	err = jsonschema.ValidateModelSchema(json)
	if err != nil {
		return
	}
	err = zjson.Unmarshal(json, &m)
	if err == nil {
		m.Raw = json
		m.readOnlyKeys = make([]string, 0)
		m.cryptKeys = make(map[string]cryptProcess, 0)
		m.afterProcess = make(map[string][]afterProcess, 0)
		m.beforeProcess = make(map[string][]beforeProcess, 0)

		m.fields = zarray.Map(m.Columns, func(_ int, c *Column) string {
			parseColumn(m, c)

			parseValidRule(c)

			parseColumnOptions(c)
			return c.Name
		})

		m.inlayFields = []string{IDKey}
		if m.Options.Timestamps {
			m.inlayFields = append(m.inlayFields, CreatedAtKey, UpdatedAtKey)
		}

		if m.Options.SoftDeletes {
			m.inlayFields = append(m.inlayFields, DeletedAtKey)
		}

		m.fullFields = append([]string{IDKey}, m.fields...)
		m.fullFields = zarray.Unique(append(m.fullFields, m.inlayFields...))
	}
	return
}
