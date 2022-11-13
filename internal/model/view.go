package model

import (
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
	columns := make([]ztype.Map, 0)
	data, ok := m.Views["lists"]

	if !ok {
		return ztype.Map{"k": 3}
	}
	fields := data.Fields
	if len(fields) == 0 {
		fields = m.fields
	}
	for _, v := range fields {
		column, ok := m.getColumn(v)
		if !ok {
			continue
		}
		columns = append(columns, ztype.Map{
			"title": column.Label,
			"key":   column.Name,
			"type":  column.Type,
		})
	}
	info := ztype.Map{
		"title":   zutil.IfVal(data.Title != "", data.Title, m.Name+"数据"),
		"columns": columns,
		"fields":  fields,
	}
	return info
}
func resolverViewInfo(m *Model) ztype.Map {
	info := ztype.Map{}
	return info
}

func resolverView(m *Model) ztype.Map {
	views := ztype.Map{
		"model": m.Name,
	}

	views["lists"] = resolverViewLists(m)

	views["info"] = resolverViewInfo(m)
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
