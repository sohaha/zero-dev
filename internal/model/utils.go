package model

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
)

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
