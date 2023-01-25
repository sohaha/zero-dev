package api

import (
	"os"

	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
)

func ParseRestApi(name string, json []byte) {
	zlog.Debug(name)
	j := zjson.ParseBytes(json)

	j.Get("routes").ForEach(func(_, value *zjson.Res) bool {
		zlog.Debug(value.Get("path"))
		SetRouter("GET", value.Get("path").String(), value.Get("handler").String())
		return true
	})

	zlog.Println("\n\n")
	GetRouter("GET", "/")
	GetRouter("GET", "/ss/id")
	GetRouter("GET", "/ss")
	GetRouter("GET", "/ss/3345")
	os.Exit(0)
}
