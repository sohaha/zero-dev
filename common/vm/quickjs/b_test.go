package quickjs_test

import (
	"testing"

	quickjs2 "github.com/buke/quickjs-go"
	"github.com/syumai/quickjs"
)

func TestT(t *testing.T) {

}

func BenchmarkQuickjs(b *testing.B) {
	runtime := quickjs.NewRuntime()
	defer runtime.Free()

	context := runtime.NewContext()
	defer context.Free()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := context.Eval("`Hello world! 2 ** 8 = ${2 ** 8}.`")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkQuickjsGo(b *testing.B) {
	runtime := quickjs2.NewRuntime()
	defer runtime.Close()

	context := runtime.NewContext()
	defer context.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := context.Eval("`Hello world! 2 ** 8 = ${2 ** 8}.`")
		if err != nil {
			b.Error(err)
		}
	}
}
