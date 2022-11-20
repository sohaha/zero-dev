package parse

import (
	"strings"

	"github.com/sohaha/zlsgo/znet"
)

type relation struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Model   string   `json:"model"`
	Foreign string   `json:"foreign"`
	Key     string   `json:"key"`
	Fields  []string `json:"fields"`
}

func GetRequestWiths(c *znet.Context, m *Modeler) (mustFields []string, hasOne map[string]*relation, hasMany map[string]*relation) {
	with, ok := c.GetQuery("with")
	if !ok || with == "" {
		return []string{}, map[string]*relation{}, map[string]*relation{}
	}

	w := strings.Split(with, ",")
	mustFields = make([]string, 0, len(w))
	hasOne = make(map[string]*relation, len(w))
	hasMany = make(map[string]*relation, len(w))

	for _, v := range w {
		r, ok := m.Relations[v]
		if !ok {
			continue
		}
		mustFields = append(mustFields, m.Table.Name+"."+r.Key)
		if r.Type == "hasMany" {
			hasMany[v] = r
		} else {
			hasOne[v] = r
		}
	}

	return
}
