package quickjs

import (
	"sync"

	"github.com/buke/quickjs-go"
	polyfill "github.com/buke/quickjs-go-polyfill"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zstring"
	// "github.com/syumai/quickjs"
)

func RunFile(path string) (string, error) {
	return "", nil
}

type JS struct {
	runtime *quickjs.Runtime
	ctx     *quickjs.Context
	pool    sync.Pool
}
type runtime struct {
	runtime quickjs.Runtime
	ctx     *quickjs.Context
}
type Value = quickjs.Value

func New() *JS {
	// runtime := quickjs.NewRuntime()

	return &JS{
		// runtime: &runtime,
		// ctx:     ctx,
		pool: sync.Pool{
			New: func() interface{} {
				r := quickjs.NewRuntime()
				ctx := r.NewContext()
				polyfill.InjectAll(ctx)
				return &runtime{
					runtime: r,
					ctx:     ctx,
				}
			},
		},
	}
}

func (j *JS) RunScript(script string, fn func(*quickjs.Value, error)) {
	// j.runtime.
	// context := j.runtime.NewContext()

	ctx := j.pool.Get().(*quickjs.Context)
	defer j.pool.Put(ctx)
	res, err := ctx.Eval(script)

	fn(&res, err)

	zlog.Debug(err)
	zlog.Debug(res)
	res.Free()
}

func RunLocalFile[T any](j *JS, file string, fn func(*Value) (T, error)) (v T, err error) {
	var bytes []byte
	bytes, err = zfile.ReadFile(file)
	if err != nil {
		return
	}
	r := j.pool.Get().(*runtime)
	defer j.pool.Put(r)

	script := zstring.Bytes2String(bytes)
	script = `let exports = {};const main = (async()=>{` + script + `})();`

	var ret quickjs.Value
	ret, err = r.ctx.Eval(script)
	if err != nil {
		return
	}
	defer ret.Free()
	err = r.runtime.ExecuteAllPendingJobs()
	if err != nil {
		return
	}

	var exp quickjs.Value
	exp, err = r.ctx.Eval("exports")
	if err != nil {
		return
	}
	defer exp.Free()

	v, err = fn(&exp)

	return
}

func RunLocalScript[T any](j *JS, script string, fn func(*Value) (T, error)) (v T, err error) {
	script = `let exports = {};const main = (async()=>{` + script + `})();`
	v, err = RunScript(j, script, fn)
	return
}

func RunScript[T any](j *JS, script string, fn func(*Value) (T, error)) (v T, err error) {
	r := j.pool.Get().(*runtime)
	defer j.pool.Put(r)

	var ret quickjs.Value
	ret, err = r.ctx.Eval(script)
	if err != nil {
		e, _ := err.(*quickjs.Error)
		zlog.Debug(e.Cause)
		zlog.Debug(e.Stack)
		return
	}
	defer ret.Free()
	// }

	// r.runtime.ExecutePendingJob()
	// zlog.Debug(r.runtime.IsJobPending())
	// zlog.Debug(r.runtime.IsLoopJobPending())
	err = r.runtime.ExecuteAllPendingJobs()
	if err != nil {
		return
	}
	// time.Sleep(time.Second * 1)
	// rt.ExecuteAllPendingJobs()
	return fn(&ret)
	// zlog.Debug(err)
	// zlog.Debug(res)

}

func (j *JS) RunFile(path string) {
	context := j.runtime.NewContext()

	zlog.Debug(context, path)
}
