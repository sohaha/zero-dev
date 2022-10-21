package main

import (
	"testing"
)

func TestXxx(t *testing.T) {
	m := map[string]interface{}{"a": 123, "b": 456}
	t.Log(m)

	tm(m)

	t.Log(m)
}

func tm(m map[string]interface{}) {
	delete(m, "a")
}
