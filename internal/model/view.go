package model

import (
	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zutil"
)

type View struct {
	Title  string   `json:"title"`
	Fields []string `json:"fields"`
	// Filters 过滤器
	Filters []interface{} `json:"filters"`
}

func resolverViewLists(m *Model) ztype.Map {
	columns := make(map[string]ztype.Map, 0)
	data, ok := m.Views["lists"]

	fields := []string{}
	if ok {
		fields = data.Fields
		// return ztype.Map{}
	}

	if len(fields) == 0 {
		fields = m.GetFields()
	}

	for _, v := range fields {
		column, ok := m.getColumn(v)
		if !ok {
			continue
		}
		columns[column.Name] = ztype.Map{
			"title": column.Label,
			"type":  column.Type,
		}
	}

	info := ztype.Map{
		"title":   zutil.IfVal(data.Title != "", data.Title, m.Name+""),
		"columns": columns,
		"fields":  fields,
	}
	return info
}
func resolverViewInfo(m *Model) ztype.Map {
	info := ztype.Map{}

	data, ok := m.Views["detail"]
	columns := make(map[string]ztype.Map, 0)

	fields := []string{}
	if ok {
		fields = data.Fields
		// return ztype.Map{}
	}

	if len(fields) == 0 {
		fields = m.GetFields()
	}

	for _, v := range fields {
		column, ok := m.getColumn(v)
		if !ok {
			continue
		}

		r := zarray.Contains(m.readOnlyKeys, v)
		if !r {
			r = m.isInlayField(v)
		}
		columns[column.Name] = ztype.Map{
			"label":    column.Label,
			"type":     column.Type,
			"readonly": r,
			// "component": "NInput",
		}
	}

	info["columns"] = columns
	info["fields"] = fields
	return info
}

func resolverView(m *Model) ztype.Map {
	views := ztype.Map{
		"model": m.Name,
	}

	views["lists"] = resolverViewLists(m)

	views["detail"] = resolverViewInfo(m)
	// for k, v := range vs {
	// 	zlog.Debug(k, v)
	// }
	return views
}

// TODO dev
func fillView(m *Model) {
	// m.views = resolverView(m)
	// if m.Views == nil {
	// 	m.Views = make(map[string]*View)
	// }

	// name := m.Name
	// fields := m.fields
	// if v, ok := m.Views["lists"]; !ok {
	// 	m.Views["lists"] = &View{
	// 		Title:  name + "列表",
	// 		Fields: []string{},
	// 	}
	// } else {
	// 	if v.Title == "" {
	// 		v.Title = name + "列表"
	// 	}
	// 	if len(v.Fields) == 0 {
	// 		v.Fields = fields
	// 	}
	// }
}
