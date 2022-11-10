package model

import (
	"strings"

	"github.com/sohaha/zlsgo/znet"
)

type relation struct {
	Name    string   `json:"name"`
	Model   string   `json:"model"`
	Foreign string   `json:"foreign"`
	Key     string   `json:"key"`
	Fields  []string `json:"fields"`
}

func convertRelation(m *Model) {

}

func GetRequestWiths(c *znet.Context, m *Model) map[string]*relation {
	with, ok := c.GetQuery("with")
	if !ok || with == "" {
		return map[string]*relation{}
	}

	w := strings.Split(with, ",")
	rr := make(map[string]*relation, len(w))
	for _, v := range w {
		r, ok := m.Relations[v]
		if !ok {
			continue
		}
		rr[v] = r
	}

	return rr
}
