package parse_test

import (
	"testing"
	"zlsapp/internal/parse"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zfile"
)

func TestParseModel(t *testing.T) {
	tt := zlsgo.NewTest(t)

	m, err := testParseModel()
	tt.NoError(err)
	tt.Equal("新闻", m.Name)
	fields := m.GetFields()
	t.Log(fields)
	for _, v := range fields {
		c, ok := m.GetColumn(v)
		tt.EqualTrue(ok)
		t.Log(c.GetLabel(), c.Name, c.Type)
	}
}

var json, _ = zfile.ReadFile("../../testdata/news.model.json")

func testParseModel() (m *parse.Modeler, err error) {
	return parse.ParseModel(json)
}
