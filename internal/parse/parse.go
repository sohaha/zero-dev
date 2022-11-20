package parse

import (
	"encoding/json"
	"zlsapp/internal/parse/jsonschema"

	"github.com/zlsgo/zdb/schema"
)

func resolverColumn(m *Modeler, c *Column) {
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
func ParseModel(j []byte) (m *Modeler, err error) {
	err = jsonschema.ValidateModelSchema(j)
	if err != nil {
		return
	}
	err = json.Unmarshal(j, &m)

	if err == nil {
		m.Raw = j
		InitModel(m)
	}
	return
}
