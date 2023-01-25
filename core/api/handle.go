package api

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type HandleFunc func(HandleArgs) func(c *znet.Context) (ztype.Map, error)
type HandleArgs struct {
	Model string
}

var globalHandle = zarray.NewHashMap[string, HandleFunc]()

func InitGlobalHandle() {
	globalHandle.Set("model.find", func(a HandleArgs) func(c *znet.Context) (ztype.Map, error) {
		return func(c *znet.Context) (ztype.Map, error) {
			return nil, nil
		}
	})
}

func ParseHandle(handle string) (h func(c *znet.Context) (ztype.Map, error), err error) {
	if handle == "" {
		return nil, errors.New("handle is empty")
	}

	handleStr := strings.Split(handle, ".")
	l := len(handleStr)
	if l > 2 {
		switch {
		case handleStr[0] == "model":
			zlog.Debug("model", handleStr)
		}
	}

	// if globalHandle.Has(handle) {
	// 	return globalHandle.Get(handle), nil
	// }

	return func(c *znet.Context) (ztype.Map, error) {
		return ztype.Map{}, nil
	}, nil
}
