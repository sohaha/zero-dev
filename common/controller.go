package common

import (
	"reflect"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type Controller func(*znet.Context) (ztype.Map, error)

func (c *Controller) Invoke([]interface{}) ([]reflect.Value, error) {
	return []reflect.Value{reflect.ValueOf(c)}, nil
}
