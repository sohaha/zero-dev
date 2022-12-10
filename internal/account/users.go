package account

import (
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type Users struct {
	service.App
	Path  string
	model *parse.Modeler
}

func (u *Users) Init(z *znet.Engine) {
	u.model, _ = parse.GetModel(UsersModel)
}

func (u *Users) KeyGet(c *znet.Context) (interface{}, error) {
	id := c.GetParam("key")
	item, err := parse.FindOne(u.model, id)
	if err != nil {
		return nil, zerror.InvalidInput.Wrap(err, "Invalid id")
	}
	return item, nil
}

func (u *Users) Get(c *znet.Context) (interface{}, error) {
	page, size, err := parse.GetPages(c)
	if err != nil {
		return nil, zerror.InvalidInput.Wrap(err, "Invalid page or size")
	}

	filter := ztype.Map{}
	items, p, err := parse.Pages(u.model, page, size, filter)

	return ztype.Map{
		"items": items,
		"page":  p,
	}, err
}
