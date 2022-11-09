package model

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
)

func GetRequestFields(c *znet.Context, m *Model) []string {
	tableFields := m.columnsKeys
	tableFields = append(tableFields, IDKey)
	if m.Options.Timestamps {
		tableFields = append(tableFields, CreatedAtKey, UpdatedAtKey)
	}

	f, ok := c.GetQuery("fields")
	if !ok || f == "" {
		return tableFields
	}

	fields := strings.Split(f, ",")
	if len(fields) == 0 {
		return tableFields
	}

	return zarray.Filter(fields, func(_ int, f string) bool {
		return zarray.Contains(tableFields, f)
	})

}
