package open

import (
	"zlsapp/internal/model"

	"github.com/sohaha/zlsgo/znet"
)

func (h *Open) GetSchemaModel(c *znet.Context) {
	c.Byte(200, model.GetModelSchema())
	c.SetContentType(znet.ContentTypeJSON)
}
