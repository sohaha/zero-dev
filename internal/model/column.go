package model

import (
	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zvalid"
	"github.com/zlsgo/zdb/schema"
)

type Column struct {
	Default     interface{}     `json:"default"`
	Unique      interface{}     `json:"unique"`
	Index       interface{}     `json:"index"`
	Crypt       string          `json:"crypt"`
	Name        string          `json:"name"`
	Comment     string          `json:"comment"`
	Label       string          `json:"label"`
	Type        schema.DataType `json:"type"`
	Tag         string          `json:"tag"`
	Validations []validations   `json:"validations"`
	Options     []ColumnEnum    `json:"options"`
	Before      []string        `json:"before"`
	After       []string        `json:"after"`
	validRules  zvalid.Engine   `json:"-"`
	Size        uint64          `json:"size"`
	ReadOnly    bool            `json:"read_only"`
	Nullable    bool            `json:"nullable"`
}

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

func (m *Model) GetFields(exclude ...string) []string {
	if len(exclude) == 0 {
		return m.fields
	}

	return zarray.Filter(m.fields, func(_ int, v string) bool {
		return !zarray.Contains(exclude, v)
	})
}

func (m *Model) getColumn(name string) (*Column, bool) {
	column, ok := zarray.Find(m.Columns, func(_ int, c *Column) bool {
		return c.Name == name
	})
	if ok {
		return column, true
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

	return nil, false
}

type ColumnEnum struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
