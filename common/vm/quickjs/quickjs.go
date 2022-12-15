package quickjs

import (
	"github.com/sohaha/zlsgo/zlog"
	"github.com/syumai/quickjs"
)

func RunFile(path string) (string, error) {
	return "", nil
}

type JS struct {
	runtime quickjs.Runtime
}

func New() *JS {
	runtime := quickjs.NewRuntime()
	return &JS{
		runtime: runtime,
	}
}

func (j *JS) RunScript(script string) {
	// j.runtime.
	context := j.runtime.NewContext()

	zlog.Debug(context.Eval(script))
}

func (j *JS) RunFile(path string) {
	context := j.runtime.NewContext()

	zlog.Debug(context, path)
}
