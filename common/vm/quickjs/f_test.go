package quickjs_test

import (
	"sync"
	"testing"
	"zlsapp/common/vm/quickjs"

	"github.com/sohaha/zlsgo/ztype"
)

func TestFile(t *testing.T) {
	js := quickjs.New()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			v, err := quickjs.RunLocalFile(js, "./testdata/1.js", func(v *quickjs.Value) (ztype.Map, error) {
				return ztype.Map{
					"now":  v.Get("default").Get("now").String(),
					"rand": v.Get("default").Get("rand").Int64(),
				}, nil
			})

			t.Log(v, err)
			wg.Done()
		}()
	}

	wg.Wait()
}
