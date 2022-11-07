package common

import (
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

// type Controller func(*znet.Context) (ztype.Map, error)

// func (c *Controller) Invoke([]interface{}) ([]reflect.Value, error) {
// 	return []reflect.Value{reflect.ValueOf(c)}, nil
// }

func GetUID(c *znet.Context) string {
	id, ok := c.Value("uid", "")
	if !ok {
		return ""
	}
	return ztype.ToString(id)
}
