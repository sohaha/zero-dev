package hashmap

import (
	"testing"

	"github.com/alphadose/haxmap"
)

func TestM(t *testing.T) {
	h := haxmap.New[string, string]()

	h.Set("a", "a1")
	h.Set("b", "b2")
	h.Set("c", "c3")
	t.Log(h.Len())
	t.Log(h.Get("a"))
	t.Log(h.Get("b"))
	h.Del("a")
	t.Log(h.Len())
	t.Log(h.Get("a"))
	t.Log(h.Get("b"))
	h.Del("b")
	t.Log(h.Len())
	t.Log(h.Get("a"))
	t.Log(h.Get("b"))
	h.Del("c")
	t.Log(h.Len())
	t.Log(h.Get("a"))
	t.Log(h.Get("b"))
	t.Log(h.Get("c3"))
}
