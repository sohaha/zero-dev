package parse

import (
	"github.com/sohaha/zlsgo/znet"
)

type ModelRelation struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Model   string   `json:"model"`
	Foreign string   `json:"foreign"`
	Key     string   `json:"key"`
	Fields  []string `json:"fields"`
}

func GetRequestWiths(c *znet.Context, m *Modeler, withFilds []string) (mustFields []string, hasOne map[string]*ModelRelation, hasMany map[string]*ModelRelation) {
	if len(withFilds) == 0 {
		return []string{}, map[string]*ModelRelation{}, map[string]*ModelRelation{}
	}

	mustFields = make([]string, 0, len(withFilds))
	hasOne = make(map[string]*ModelRelation, len(withFilds))
	hasMany = make(map[string]*ModelRelation, len(withFilds))

	for _, v := range withFilds {
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
