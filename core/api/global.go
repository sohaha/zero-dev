package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
)

var methods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodPost:    {},
	http.MethodPut:     {},
	http.MethodDelete:  {},
	http.MethodPatch:   {},
	http.MethodHead:    {},
	http.MethodOptions: {},
	http.MethodConnect: {},
	http.MethodTrace:   {},
}

type Router struct {
	trees map[string]*znet.Tree
}

var t = &Router{
	trees: make(map[string]*znet.Tree),
}

var globalApis = zarray.NewHashMap[string, *Router]()

func ClearRouter() {
	globalApis = zarray.NewHashMap[string, *Router]()
	t = &Router{
		trees: make(map[string]*znet.Tree),
	}
}

func SetRouter(method string, path string, handle string) error {
	if _, ok := methods[method]; !ok {
		return errors.New(method + " is invalid method")
	}

	tree, ok := t.trees[method]
	if !ok {
		tree = znet.NewTree()
		t.trees[method] = tree
	}

	_ = tree.Add(path, func(c *znet.Context) error {
		zlog.Debug("Get")
		h, err := ParseHandle(handle)
		if err != nil {
			return err
		}

		return znet.Utils.ParseHandlerFunc(h)(c)
	}).WithValue(handle)

	// zlog.Debug(tree)
	return nil
}

func GetRouter(method string, path string) {
	tree, ok := t.trees[method]
	if !ok {
		return
	}

	nodes := tree.Find(path, false)
	for i := range nodes {
		node := nodes[i]
		zlog.Debug(path, "  == ", node.Path(), node.Path() == path)
		if node.Path() == path {
			zlog.Success("ok", path, node.Value(), "\n")

			return
		}
	}

	if len(nodes) == 0 {
		res := strings.Split(path, "/")
		p := ""
		if len(res) == 1 {
			p = res[0]
		} else {
			p = res[1]
		}
		nodes := tree.Find(p, true)
		for _, node := range nodes {
			zlog.Debug(path, "  -- ", node.Path(), node.Path() == path)
			if node.Path() != path {
				if matchParamsMap, ok := znet.Utils.URLMatchAndParse(path, node.Path()); ok {
					zlog.Success("ok", matchParamsMap, path, node.Value(), "\n")
					node.Handle()(nil)
					return
				}
			}
		}
	}
	// 找到的是 node
	zlog.Debug(path, nodes)
}
