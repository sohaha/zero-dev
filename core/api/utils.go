package api

import (
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/znet"
	"github.com/zlsgo/zdb/builder"
	"zlsapp/internal/parse"
)

func getFinalFields(m *parse.Modeler, c *znet.Context, fields []string, withFilds []string) (finalFields, tmpFields []string, with, withMany map[string]*parse.ModelRelation) {
	var mustFields []string
	mustFields, with, withMany = GetRequestWiths(c, m, withFilds)

	if len(fields) == 0 {
		fields = GetRequestFields(c, m)
	}

	finalFields = zarray.Unique(append(fields, mustFields...))
	_, tmpFields = zarray.Diff(fields, mustFields)

	return
}

func GetRequestFields(c *znet.Context, m *parse.Modeler) []string {
	tableFields := make([]string, 0, len(m.Fields)+1)
	tableFields = append(tableFields, parse.IDKey)
	tableFields = append(tableFields, m.Fields...)

	if m.Options.Timestamps {
		tableFields = append(tableFields, parse.CreatedAtKey, parse.UpdatedAtKey)
	}

	if m.Options.CreatedBy {
		tableFields = append(tableFields, parse.CreatedByKey)
	}

	if f, ok := c.GetQuery("fields"); ok && f != "" {
		fields := strings.Split(f, ",")
		if len(fields) > 0 {
			tableFields = zarray.Filter(fields, func(_ int, f string) bool {
				return zarray.Contains(tableFields, f)
			})
		}
	}

	return zarray.Map(tableFields, func(_ int, f string) string {
		return m.Table.Name + "." + f
	})
}

func GetRequestWiths(c *znet.Context, m *parse.Modeler, withFilds []string) (mustFields []string, hasOne map[string]*parse.ModelRelation, hasMany map[string]*parse.ModelRelation) {
	if len(withFilds) == 0 {
		return []string{}, map[string]*parse.ModelRelation{}, map[string]*parse.ModelRelation{}
	}

	mustFields = make([]string, 0, len(withFilds))
	hasOne = make(map[string]*parse.ModelRelation, len(withFilds))
	hasMany = make(map[string]*parse.ModelRelation, len(withFilds))

	for _, v := range withFilds {
		r, ok := m.Relations[v]
		if !ok {
			continue
		}
		mustFields = append(mustFields, m.Table.Name+"."+r.Key)
		if r.Type == "hasMany" {
			hasMany[v] = r
		} else {
			r.Join = builder.LeftJoin
			hasOne[v] = r
		}
	}

	return
}
