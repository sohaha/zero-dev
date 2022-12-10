package parse

import (
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
	Validations []Validations   `json:"validations"`
	Options     []ColumnEnum    `json:"options"`
	Before      []string        `json:"before"`
	After       []string        `json:"after"`
	validRules  zvalid.Engine   `json:"-"`
	Size        uint64          `json:"size"`
	ReadOnly    bool            `json:"readonly"`
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
