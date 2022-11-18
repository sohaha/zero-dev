package model_test

import (
	"testing"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
)

func TestModel_ActionUpdate(t *testing.T) {
	tt := zlsgo.NewTest(t)
	m, err := getModel(true)
	tt.NoError(err)

	id, err := m.ActionCreate(ztype.Map{
		"title":    "test",
		"category": 1,
		"content":  zstring.Rand(100),
	})

	t.Log(id, err)

	tt.NoError(err)
}
