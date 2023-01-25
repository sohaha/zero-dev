package parse

import (
	"github.com/sohaha/zlsgo/zarray"
)

func GetViewFields(m *Modeler, view string) []string {
	v, ok := m.Views[view]
	if !ok {
		return []string{}
	}

	fields := v.Fields
	if len(fields) == 0 {
		fields = m.GetFields()
	}
	return zarray.Unique(append(fields, IDKey))
}
