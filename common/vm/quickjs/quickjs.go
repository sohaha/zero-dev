package quickjs

import (
	"sync"
	"time"

	"github.com/buke/quickjs-go"
	polyfill "github.com/buke/quickjs-go-polyfill"
	"github.com/sohaha/zlsgo/zlog"
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

func RunScript[T any](j *JS, script string, fn func(*Value, error) (T, error)) (T, error) {
	// j.runtime.
	// context := j.runtime.NewContext()
	r := j.pool.Get().(*runtime)
	defer j.pool.Put(r)
	res, err := r.ctx.Eval(script)
	// if err == nil {
	defer res.Free()
	// }

	// r.runtime.ExecutePendingJob()
	r.runtime.ExecuteAllPendingJobs()
	time.Sleep(time.Second * 1)
	// rt.ExecuteAllPendingJobs()
	return fn(&res, err)
	// zlog.Debug(err)
	// zlog.Debug(res)

}

func (j *JS) RunFile(path string) {
	context := j.runtime.NewContext()

	zlog.Debug(context, path)
}
