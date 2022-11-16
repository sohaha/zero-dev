package yaegi

import (
	"sync"
	"testing"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func TestYaegi(t *testing.T) {
	i := interp.New(interp.Options{})

	_ = i.Use(stdlib.Symbols)

	_, err := i.Eval(`import "fmt"`)
	if err != nil {
		panic(err)
	}

	_, err = i.Eval(`fmt.Println("Hello Yaegi")`)
	if err != nil {
		panic(err)
	}
}

//go:generate go install github.com/traefik/yaegi/cmd/yaegi
//go:generate yaegi extract github.com/sohaha/zlsgo/zlog
func TestFile(t *testing.T) {
	tt := zlsgo.NewTest(t)

	// f,err := zfile.ReadFile("./plugin/fib.go")
	intp := interp.New(interp.Options{})
	_ = intp.Use(stdlib.Symbols)
	_ = intp.Use(Symbols)

	intp.ImportUsed()

	// intp.FileSet().AddFile("plugin/fib.go", -1, zfile.ReadFile("./plugin/fib.go"))
	_, err := intp.EvalPath(zfile.RealPath("./plugin"))
	tt.NoError(err)
	// _, _ = intp.EvalPath(zfile.RealPath("./plugin/fib.go"))
	v, _ := intp.Eval("plugin.Fib")
	fu := v.Interface().(func(int) int)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v int) {
			t.Log("Fib("+ztype.ToString(v)+") =", fu(v))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
