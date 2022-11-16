package main

import (
	"reflect"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var Symbols = map[string]map[string]reflect.Value{}

func main() {
	intp := interp.New(interp.Options{})
	_ = intp.Use(stdlib.Symbols)
	_ = intp.Use(Symbols)
	_, err := intp.EvalPath(zfile.RealPath("./plugin"))
	zlog.Debug(err)
}
