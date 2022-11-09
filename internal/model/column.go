package model

import (
	"sync"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zvalid"
)

type Column struct {
	once        sync.Once
	Name        string        `json:"name"`
	Comment     string        `json:"comment"`
	Type        string        `json:"type"`
	Size        uint64        `json:"size"`
	Tag         string        `json:"tag"`
	Nullable    bool          `json:"nullable"`
	Label       string        `json:"label"`
	Enum        interface{}   `json:"enum"`
	Default     interface{}   `json:"default"`
	Unique      interface{}   `json:"unique"`
	Index       interface{}   `json:"index"`
	Validations []validations `json:"validations"`
	validRules  zvalid.Engine `json:"-"`
	ReadOnly    bool          `json:"read_only"` // 是否创建之后不允许更改
	Side        bool          `json:"side"`
	// 加密字段
	Crypt  string   `json:"crypt"`
	Before []string `json:"before"`
	After  []string `json:"after"`
}

func (c *Column) GetValidations() zvalid.Engine {
	c.once.Do(func() {
		name := c.Name
		label := c.Label
		if label == "" {
			label = name
		}
		c.validRules = parseValidRule(label, c.Validations, c.Size)
	})

	return c.validRules
}

func (m *Model) GetFields(exclude ...string) []string {
	if len(exclude) == 0 {
		return m.columnsKeys
	}

	return zarray.Filter(m.columnsKeys, func(_ int, v string) bool {
		return !zarray.Contains(exclude, v)
	})
}