package main

import (
	"encoding/json"
	"testing"

	"github.com/sohaha/zlsgo/zjson"
)

func TestXX(t *testing.T) {
	t.Log(zjson.Valid("dd"))
	t.Log(zjson.Valid("33"))
	t.Log(json.Valid([]byte("33")))
}
