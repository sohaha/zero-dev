package lua

import (
	"sync"
	"testing"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/vlorc/lua-vm/base"
	"github.com/vlorc/lua-vm/pool"
	lua "github.com/yuin/gopher-lua"
)

func TestLua(t *testing.T) {

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			L := luaPool.Get()
			defer luaPool.Put(L)
			err := L.DoFile(zfile.RealPath("plugin/demo.lua"))
			if err != nil {
				t.Log(err)
			}
			c := lua.P{
				Fn:      L.GetGlobal("fib2"),
				NRet:    1,
				Protect: true,
			}
			err = L.CallByParam(c, lua.LNumber(10))
			ret := L.Get(-1)
			L.Pop(1)
			zlog.Warn()
			zlog.Debug(err, ret)
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
