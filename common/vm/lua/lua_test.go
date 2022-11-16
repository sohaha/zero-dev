package lua

import (
	"sync"
	"testing"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/vlorc/lua-vm/base"
	"github.com/vlorc/lua-vm/pool"
)

func TestLua(t *testing.T) {
	state := luaPool.Get()
	defer luaPool.Put(state)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := state.DoFile(zfile.RealPath("plugin/demo.lua"))
			if err != nil {
				t.Log(err)

			}
		}()
	}
	wg.Wait()
}
func TestLua2(t *testing.T) {

	p := pool.NewLuaPool().Preload(
		pool.Library(),

		// pool.Module("net.tcp", tcp.NewTCPFactory(driver.DirectDriver{})),
		pool.Module("buffer", base.BufferFactory{}),
		pool.Module("time", base.TimeFactory{}),
	)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {

			defer wg.Done()
			m, err := p.ModuleFile(zfile.RealPath("plugin/demo.lua"))
			if nil != err {
				println("error: ", err.Error())
			}
			// var fn func(int) int
			var count int
			t.Log(m.To(&count, "count"))
			// ok := m.Method("fib", &fn)
			// t.Log(ok)
			// t.Log(fn)
		}()
	}
	wg.Wait()
}
