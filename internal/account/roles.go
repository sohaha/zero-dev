package account

import (
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
)

type Roles struct {
	service.App
	Path  string
	model *parse.Modeler
}

func (r *Roles) Init(z *znet.Engine) {
	r.model, _ = parse.GetModel(RolesModel)
}

func (r *Roles) KeyGet(c *znet.Context) (any, error) {
	id := c.GetParam("key")
	item, err := parse.FindOne(r.model, id)
	if err != nil {
		return nil, zerror.InvalidInput.Wrap(err, "Invalid id")
	}
	return item, nil
}

func (r *Roles) Get(c *znet.Context) (any, error) {
	filter := ztype.Map{}
	items, err := parse.Find(r.model, filter, func(so *parse.StorageOptions) error {
		so.Fields = []string{parse.IDKey, "name", "key"}
		return nil
	})
	if err != nil && err == zdb.ErrNotFound {
		err = nil
	}
	return items, err
}

func (r *Roles) KeyDelete(c *znet.Context) (any, error) {
	id := c.GetParam("key")
	_, err := parse.Delete(r.model, id)
	if err != nil {
		return nil, zerror.InvalidInput.Wrap(err, "Invalid id")
	}
	return nil, nil
}
