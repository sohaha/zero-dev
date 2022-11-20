package parse

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
)

func getFinalFields(m *Modeler, c *znet.Context, fields []string) (finalFields, tmpFields []string, with, withMany map[string]*relation) {
	var mustFields []string
	mustFields, with, withMany = GetRequestWiths(c, m)

	if len(fields) == 0 {
		fields = GetRequestFields(c, m)
	}

	finalFields = zarray.Unique(append(fields, mustFields...))
	_, tmpFields = zarray.Diff(fields, mustFields)
	return
}

func GetViewFields(m *Modeler, view string) []string {
	v, ok := m.Views[view]
	if !ok {
		return []string{}
	}

	fields := v.Fields
	if len(fields) == 0 {
		fields = m.GetFields()
	}
	return zarray.Unique(append(fields, IDKey))
}

func GetRequestFields(c *znet.Context, m *Modeler) []string {
	tableFields := make([]string, 0, len(m.fields)+1)
	tableFields = append(tableFields, IDKey)
	tableFields = append(tableFields, m.fields...)

	if m.Options.Timestamps {
		tableFields = append(tableFields, CreatedAtKey, UpdatedAtKey)
	}

	if f, ok := c.GetQuery("fields"); ok && f != "" {
		fields := strings.Split(f, ",")
		if len(fields) > 0 {
			tableFields = zarray.Filter(fields, func(_ int, f string) bool {
				return zarray.Contains(tableFields, f)
			})
		}
	}

	return tableFields
}
