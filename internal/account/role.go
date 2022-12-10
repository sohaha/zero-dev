package account

import (
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type Role struct {
	service.App
	Path  string
	model *parse.Modeler
}

func (r *Role) Init(z *znet.Engine) {
	r.model, _ = parse.GetModel(RolesModel)
}

func (r *Role) KeyGet(c *znet.Context) (interface{}, error) {
	id := c.GetParam("key")
	item, err := parse.FindOne(r.model, id)
	if err != nil {
		return nil, zerror.InvalidInput.Wrap(err, "Invalid id")
	}
	return item, nil
}

func (r *Role) Get(c *znet.Context) (interface{}, error) {
	page, size, err := parse.GetPages(c)
	if err != nil {
		return nil, zerror.InvalidInput.Wrap(err, "Invalid page or size")
	}

	filter := ztype.Map{}
	items, p, err := parse.Pages(r.model, page, size, filter)

	return ztype.Map{
		"items": items,
		"page":  p,
	}, err
}
