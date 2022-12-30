package common

import (
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

func GetUID(c *znet.Context) string {
	id, ok := c.Value("uid", "")
	if !ok {
		return ""
	}
	return ztype.ToString(id)
}
