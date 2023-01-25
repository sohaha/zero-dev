package parse

import (
	"github.com/zlsgo/zdb/builder"
)

type ModelRelation struct {
	Name    string             `json:"name"`
	Type    string             `json:"type"`
	Join    builder.JoinOption `json:"-"`
	Model   string             `json:"model"`
	Foreign string             `json:"foreign"`
	Key     string             `json:"key"`
	Fields  []string           `json:"Fields"`
}
