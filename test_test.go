package main

import (
	"strings"
	"testing"
	"time"

	"github.com/sohaha/zlsgo/ztime"
)

func TestXxx(t *testing.T) {
	validDate := "2015.06.08-2022.06.08"
	v := strings.Split(validDate, "-")
	validDate = v[len(v)-1]
	expirationDate, _ := ztime.Parse(validDate, "Y.m.d")
	t.Log(expirationDate.Before(time.Now()))
	t.Log(validDate)
	// expirationDate 小于当前时间
	t.Log(expirationDate.Before(time.Now()))
	// m := map[string]interface{}{"a": 123, "b": 456}
	// t.Log(m)

	// tm(m)

	// t.Log(m)
}

func tm(m map[string]interface{}) {
	delete(m, "a")
}
