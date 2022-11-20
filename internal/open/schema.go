package open

import (
	"github.com/sohaha/zlsgo/znet"
)

func (h *Open) GetSchemaModel(c *znet.Context) {
	// c.Byte(200, parse.GetModelSchema())
	c.SetContentType(znet.ContentTypeJSON)
}
