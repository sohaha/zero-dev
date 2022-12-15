package quickjs_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
	"zlsapp/common/vm/quickjs"
	// "github.com/syumai/quickjs"
)

// func check(err error) {
// 	if err != nil {
// 		var evalErr *quickjs.Error
// 		if errors.As(err, &evalErr) {
// 			fmt.Println(evalErr.Cause)
// 			fmt.Println(evalErr.Stack)
// 		}
// 		panic(err)
// 	}
// }

func TestQuickjs(t *testing.T) {
	js := quickjs.New()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; utf-8")
		time.Sleep(time.Second * 2)
		_, _ = w.Write([]byte(`{"status": true}`))
	}))
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		// id := i
		wg.Add(1)
		go func() {
			a, err := quickjs.RunScript(js, fmt.Sprintf(`
console.log(+new Date());
		var b = async()=>{
		return 	await fetch('%s', {Method: 'GET'}).then(response => response.json()).then(data => {
			console.log(data.status);
			return 12;
		})
		}

		b()
	`, srv.URL), func(v *quickjs.Value) (string, error) {
				return v.String(), nil
			})
			t.Log(a, err)
			// js.RunScript(`("!` + ztype.ToString(i) + `")`)
			// js.RunScript(`console.log("!` + ztype.ToString(i) + `")`)
			// js.RunScript(`console.error("Error Hello world!` + ztype.ToString(i) + `")`)
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(time.Second * 3)
	return
	// runtime := quickjs.NewRuntime()
	// defer runtime.Free()

	// context := runtime.NewContext()
	// defer context.Free()

	// globals := context.Globals()

	// // Test evaluating template strings.

	// result, err := context.Eval("`Hello world! 2 ** 8 = ${2 ** 8}.`")
	// check(err)
	// defer result.Free()

	// fmt.Println(result.String())
	// fmt.Println()

	// // Test evaluating numeric expressions.

	// result, err = context.Eval(`1 + 2 * 100 - 3 + Math.sin(10)`)
	// check(err)
	// defer result.Free()

	// fmt.Println(result.Int64())
	// fmt.Println()

	// // Test evaluating big integer expressions.

	// result, err = context.Eval(`128n ** 16n`)
	// check(err)
	// defer result.Free()

	// fmt.Println(result.BigInt())
	// fmt.Println()

	// // Test evaluating big decimal expressions.

	// result, err = context.Eval(`128l ** 12l`)
	// check(err)
	// defer result.Free()

	// fmt.Println(result.BigFloat())
	// fmt.Println()

	// // Test evaluating boolean expressions.

	// result, err = context.Eval(`false && true`)
	// check(err)
	// defer result.Free()

	// fmt.Println(result.Bool())
	// fmt.Println()

	// // Test setting and calling functions.

	// A := func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	// 	fmt.Println("A got called!", args[0].Int64())
	// 	o := ctx.Object()
	// 	o.Set("ak", ctx.String("av"))
	// 	return o
	// }

	// B := func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	// 	fmt.Println("B got called!", args)
	// 	return ctx.Null()
	// }

	// globals.Set("A", context.Function(A))
	// globals.Set("B", context.Function(B))
	// globals.Set("log", context.Function(func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
	// 	t.Log(args[1].IsObject())
	// 	t.Log(args[1].PropertyNames())
	// 	t.Log(args[1].Get("ak").String())
	// 	return ctx.Null()
	// }))

	// _, err = context.Eval(`for (let i = 0; i < 3; i++) { if (i % 2 === 0) log("a--",A("a1")); else B(); }`)
	// check(err)

	// fmt.Println()

	// // // Test setting global variables.

	// // result, err = context.Eval(`HELLO = "world"; TEST = false;`)
	// // check(err)
	// // t.Log(result)

	// // names, err := globals.PropertyNames()
	// // check(err)

	// // fmt.Println("Globals:")
	// // for _, name := range names {
	// // 	val := globals.GetByAtom(name.Atom)
	// // 	defer val.Free()

	// // 	fmt.Printf("'%s': %s\n", name, val)
	// // }
	// // fmt.Println()

	// // Test evaluating arbitrary expressions from flag arguments.

	// flag.Parse()
	// if flag.NArg() == 0 {
	// 	return
	// }

	// result, err = context.Eval(strings.Join(flag.Args(), " "))
	// check(err)
	// defer result.Free()

	// if result.IsObject() {
	// 	names, err := result.PropertyNames()
	// 	check(err)

	// 	fmt.Println("Object:")
	// 	for _, name := range names {
	// 		val := result.GetByAtom(name.Atom)
	// 		defer val.Free()

	// 		fmt.Printf("'%s': %s\n", name, val)
	// 	}
	// } else {
	// 	fmt.Println(result.String())
	// }

}

func TestForFile(t *testing.T) {

}
