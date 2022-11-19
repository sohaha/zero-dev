package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/ztype"
	"golang.org/x/exp/constraints"
)

func TestXX(t *testing.T) {
	t.Log(zjson.Valid("dd"))
	t.Log(zjson.Valid("33"))
	t.Log(json.Valid([]byte("33")))
}

type tf interface {
	ztype.Map | constraints.Integer
}

func find[T tf](f T) {
	var v interface{} = f
	switch val := v.(type) {
	case ztype.Map:
		fmt.Sprintln(val.Get("a"))
	}
}

func findMap(f ztype.Map) {
	fmt.Sprintln(f.Get("a"))
}

func BenchmarkXxx1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		find(ztype.Map{"a": 1})
	}
}

func BenchmarkXxxMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		find(ztype.Map{"a": 1})
	}
}
