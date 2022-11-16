package tengo

import (
	"context"
	"fmt"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/sohaha/zlsgo/zfile"
)

func TestTengo(t *testing.T) {
	// create a new Script instance
	script := tengo.NewScript([]byte(
		`each := func(seq, fn) {
    for x in seq { fn(x) }
}

sum := 0
mul := 1
each([a, b, c, d], func(x) {
    sum += x
    mul *= x
})`))

	// set values
	_ = script.Add("a", 1)
	_ = script.Add("b", 9)
	_ = script.Add("c", 8)
	_ = script.Add("d", 4)

	// run the script
	compiled, err := script.RunContext(context.Background())
	if err != nil {
		panic(err)
	}

	// retrieve values
	sum := compiled.Get("sum")
	mul := compiled.Get("mul")
	fmt.Println(sum, mul) // "22 288"
}

func TestFile(t *testing.T) {
	modules := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	_ = modules
	b, _ := zfile.ReadFile("./test.tengo")
	script := tengo.NewScript(b)

	moduleMap := stdlib.GetModuleMap(stdlib.AllModuleNames()...)
	script.SetImports(moduleMap)
	// script.SetImports(modules...)
	compiled, err := script.RunContext(context.Background())
	t.Log(err)
	v := compiled.Get("main")
	t.Log(v.Value())
}
