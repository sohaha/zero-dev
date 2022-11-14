package model

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
)

func getFinalFields(m *Model, c *znet.Context, fields []string) (finalFields, tmpFields []string, quote bool, with, withMany map[string]*relation) {
	var mustFields []string
	mustFields, with, withMany = GetRequestWiths(c, m)
	hasWith := len(with) > 0
	hasWithMany := len(withMany) > 0
	quote = hasWith || hasWithMany

	if len(fields) == 0 {
		fields = GetRequestFields(c, m, quote)
	}
	fields = append(fields, IDKey)
	finalFields = zarray.Unique(append(fields, mustFields...))
	_, tmpFields = zarray.Diff(fields, mustFields)
	return
}

func GetViewFields(m *Model, view string) []string {
	v, ok := m.Views[view]
	if !ok {
		return []string{}
	}
	return v.Fields
}

func GetRequestFields(c *znet.Context, m *Model, quote bool) []string {
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

	if !quote {
		return tableFields
	}

	return zarray.Map(tableFields, func(_ int, v string) string {
		return m.Table.Name + "." + v
	})

}
