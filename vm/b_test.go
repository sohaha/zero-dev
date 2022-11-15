package vm

import (
	"context"
	"fmt"
	"testing"

	"github.com/dop251/goja"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tx7do/go-js"
	"github.com/wapc/wapc-go"
	wazeroe "github.com/wapc/wapc-go/engines/wazero"
	// "github.com/wapc/wapc-go/engines/wazero"
)

func BenchmarkWa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Wa()
	}
}
func BenchmarkJs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Js()
	}
}
func TestWa(t *testing.T) {
	var err error
	err = Wa()
	t.Log(err)

	err = Js()
	t.Log(err)
	err = Js2()
	t.Log(err)
	err = Waz()
	t.Log(err)
}

func Wa() error {
	ctx := context.TODO()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	addWasm, err := zfile.ReadFile("./go/add.wasm")
	// fmt.Println(err, len(addWasm))
	_ = err
	mod, err := r.InstantiateModuleFromBinary(ctx, addWasm)
	if err != nil {
		// fmt.Println(err)
		return err
	}
	x, y := uint64(1), uint64(2)

	add := mod.ExportedFunction("add")
	results, err := add.Call(ctx, x, y)
	if err != nil {

		return err
	}

	_ = results[0]
	// fmt.Printf("%d + %d = %d\n", x, y, results[0])

	return nil
}

func Js() error {
	var script = `
	function add(x,y) {
		return x + y
	}
	`
	vm := goja.New()

	_, err := vm.RunString(script)
	if err != nil {
		return err
	}

	var fn func(i, j uint64) uint
	err = vm.ExportTo(vm.Get("add"), &fn)
	if err != nil {
		return err
	}

	x, y := uint64(1), uint64(2)
	results := fn(x, y)
	_ = results
	// fmt.Printf("%d + %d = %d\n", x, y, results)
	return nil
}

func Waz() error {

	ctx := context.TODO()
	engine := wazeroe.Engine()

	add, _ := zfile.ReadFile("./go/add.wasm")
	module, err := engine.New(ctx, func(ctx context.Context, binding, namespace, operation string, payload []byte) ([]byte, error) {
		return payload, nil
	}, add, &wapc.ModuleConfig{
		Logger: wapc.PrintlnLogger,
		// Stdout: os.Stdout,
		// Stderr: os.Stderr,
	})
	if err != nil {
		return err
	}
	defer module.Close(ctx)

	instance, err := module.Instantiate(ctx)
	if err != nil {
		return err
	}
	defer instance.Close(ctx)

	result, err := instance.Invoke(ctx, "add", []byte("name"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result))
	return nil
}

func Js2() error {
	exe := js.NewVirtualMachine()
	var script = `
	function add(x,y) {
		return x + y
	}
	`
	err := exe.LoadString(script)
	if err != nil {
		return err
	}

	err = exe.Execute()
	if err != nil {
		return err
	}

	var fn func(i, j uint64) uint
	err = exe.GetFunction("add", &fn)
	if err != nil {
		return err
	}

	r := fn(1, 2)
	_ = r
	return nil
}
