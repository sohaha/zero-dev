package parse

import (
	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zvalid"
	"github.com/zlsgo/zdb/schema"
)

func (c *Column) GetValidations() zvalid.Engine {
	return c.validRules
}

func (c *Column) GetLabel() string {
	label := c.Label
	if label == "" {
		label = c.Name
	}
	return label
}

func (m *Model) isInlayField(field string) bool {
	if field == IDKey {
		return true
	}
	if !m.Options.Timestamps {
		return false
	}

	return field == CreatedAtKey || field == UpdatedAtKey
}

func (m *Model) GetFields(exclude ...string) []string {
	f := m.fullFields
	if len(exclude) == 0 {
		return f
	}

	return zarray.Filter(f, func(_ int, v string) bool {
		return !zarray.Contains(exclude, v)
	})
}

func (m *Model) GetColumn(name string) (*Column, bool) {
	column, ok := zarray.Find(m.Columns, func(_ int, c *Column) bool {
		return c.Name == name
	})
	if ok {
		return column, true
	}

	if name == IDKey {
		return &Column{
			Name:     IDKey,
			Type:     schema.Int,
			Nullable: false,
			Label:    "ID",
			ReadOnly: true,
		}, true

	}
	if m.Options.Timestamps {
		switch name {
		case CreatedAtKey:
			return &Column{
				Name:  name,
				Type:  schema.Time,
				Label: "创建时间"}, true
		case UpdatedAtKey:
			return &Column{
				Name:  name,
				Type:  schema.Time,
				Label: "更新时间"}, true
		}
	}

	if m.Options.SoftDeletes {
		if name == DeletedAtKey {
			return &Column{
				Name:  name,
				Type:  schema.Time,
				Label: "删除时间"}, true
		}
	}

	return nil, false
}
