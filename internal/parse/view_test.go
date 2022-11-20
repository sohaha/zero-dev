package parse_test

import (
	"testing"

	"github.com/sohaha/zlsgo"
)

func TestModelView(t *testing.T) {
	tt := zlsgo.NewTest(t)

	m, err := initModel(true)
	tt.NoError(err)

	t.Log(m.GetView().Get("detail"))

}
