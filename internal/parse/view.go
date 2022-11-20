package parse

import (
	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zutil"
)

type View struct {
	Title   string        `json:"title"`
	Fields  []string      `json:"fields"`
	Filters []interface{} `json:"filters"`
}

func resolverViewLists(m *Modeler) ztype.Map {
	columns := make(map[string]ztype.Map, 0)
	data, ok := m.Views["lists"]
	if !ok {
		data = &View{}
	}

	fields := []string{IDKey}
	if ok {
		fields = append(fields, data.Fields...)
	}

	if len(fields) == 1 {
		fields = append(fields, m.GetFields()...)
	}

	fields = zarray.Unique(fields)

	for _, v := range fields {
		column, ok := m.GetColumn(v)
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

func resolverViewInfo(m *Modeler) ztype.Map {
	info := ztype.Map{}

	data, ok := m.Views["detail"]
	if !ok {
		data = &View{}
	}

	columns := make(map[string]ztype.Map, 0)

	fields := []string{IDKey}
	if ok {
		fields = append(fields, data.Fields...)
	}

	if len(fields) == 1 {
		fields = append(fields, m.GetFields()...)
	}

	fields = zarray.Unique(fields)
	for _, v := range fields {
		column, ok := m.GetColumn(v)
		if !ok {
			continue
		}
		columns[column.Name] = ztype.Map{
			"label":    column.Label,
			"type":     column.Type,
			"readonly": column.ReadOnly,
			"disabled": m.isInlayField(v),
		}
	}

	info["columns"] = columns
	info["fields"] = fields
	return info
}

func resolverView(m *Modeler) {
	m.views = ztype.Map{}

	m.views["lists"] = resolverViewLists(m)

	m.views["detail"] = resolverViewInfo(m)
}

func (m *Modeler) GetView() ztype.Map {
	return m.views
}
