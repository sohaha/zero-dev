package model

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
)

func getFinalFields(m *Model, c *znet.Context) (finalFields, tmpFields []string, quote bool, with, withMany map[string]*relation) {
	var mustFields []string
	mustFields, with, withMany = GetRequestWiths(c, m)
	hasWith := len(with) > 0
	hasWithMany := len(withMany) > 0
	quote = hasWith || hasWithMany
	fields := GetRequestFields(c, m, quote)
	finalFields = zarray.Unique(append(fields, mustFields...))
	_, tmpFields = zarray.Diff(fields, mustFields)
	return
}

func GetRequestFields(c *znet.Context, m *Model, quote bool) []string {
	tableFields := m.columnsKeys
	tableFields = append(tableFields, IDKey)
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
