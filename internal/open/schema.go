package open

import (
	"zlsapp/internal/parse/jsonschema"

	"github.com/sohaha/zlsgo/znet"
)

func (h *Open) GetSchemaModel(c *znet.Context) {
	c.Byte(200, jsonschema.GetModelSchema())
	c.SetContentType(znet.ContentTypeJSON)
}
