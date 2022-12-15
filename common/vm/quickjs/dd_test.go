package quickjs_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/buke/quickjs-go"
	polyfill "github.com/buke/quickjs-go-polyfill"
)

func TestXxx(t *testing.T) {
	// Create a new runtime
	rt := quickjs.NewRuntime()
	// defer rt.Close()

	// Create a new context
	ctx := rt.NewContext()
	defer ctx.Close()

	ctx.Globals().Set("hello", ctx.Function(func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		t.Log("hello", args)
		return ctx.String("Hello " + args[0].String())
	}))

	// Inject polyfills to the context
	polyfill.InjectAll(ctx)

	now := time.Now()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; utf-8")
		t.Log("g")
		time.Sleep(time.Second * 2)
		t.Log(time.Since(now).String())
		_, _ = w.Write([]byte(`{"status": true}`))
	}))

	ret, err := ctx.Compile(fmt.Sprintf(`
	var res = 1;
	 console.log('res:', res);
	 console.log(JSON.stringify({"s":11}));
	const myFun = async () => {
		hello(1)

    return new Promise(resolve => {
		hello(3)

    //   setTimeout(() => {
		hello(2)
        const data = [
          { id: 1, name: '1', age: 11 },
          { id: 2, name: 'xiaohong', age: 22 },
          { id: 3, name: 'xiaogang', age: 33 },
        ];
        resolve(data);
    //   }, 1000);
    });
  }
  const myFun2 = async () => {
  res = await 	fetch('%s', {Method: 'GET'}).then(response => response.json()).then(data => {
			console.log(data.status);
			return 333
		});
//   res = await myFun();
  console.log('res:', res);
}
myFun2();
	`, srv.URL))
	t.Log(err)
	_ = ret
	// defer ret.Free()

	rt2 := quickjs.NewRuntime()
	// defer rt2.Close()

	// Create a new context
	ctx2 := rt2.NewContext()
	polyfill.InjectAll(ctx2)
	// defer ctx2.Close()

	//Eval bytecode
	result, err := ctx2.EvalBytecode(ret)
	t.Log(err)
	t.Log(result)
	_ = result
	// fmt.Println(result.String())
	// t.Log(ret)
	// time.Sleep(time.Millisecond * 100)
	// c, _ := rt.ExecutePendingJob()
	// t.Log(c)
	t.Log(rt.IsLoopJobPending(), rt.IsJobPending())
	t.Log(rt2.IsLoopJobPending(), rt2.IsJobPending())
	// time.Sleep(time.Second * 1)
	// rt.ExecuteAllPendingJobs()
	rt2.ExecuteAllPendingJobs()
	// t.Log(rt.IsLoopJobPending(), rt.IsJobPending())

	time.Sleep(time.Second * 3)
	// t.Log(err)
	t.Log(ctx2.Eval(`res`))
	// time.Sleep(time.Second * 1)
	// ret.Free()
	// t.Log(11)
	// t.Log(22)
	// t.Log(ctx.Eval(`res`))
	// t.Log(rt.ExecuteAllPendingJobs())
}
