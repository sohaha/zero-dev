package grbac

import (
	"zlsapp/grbac/meta"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/conf"
)

type FileLoader struct {
	path string
}

func NewFileLoader(file string) (*FileLoader, error) {
	loader := &FileLoader{
		path: zfile.RealPath(file),
	}
	_, err := loader.Load()
	if err != nil {
		return nil, err
	}
	return loader, nil
}

func (loader *FileLoader) Load() (rules meta.Rules, err error) {
	c := conf.New(loader.path)

	err = c.Read()
	if err != nil {
		return nil, err
	}

	for _, v := range c.GetAll() {
		m := ztype.ToMap(v)
		rule := &meta.Rule{
			Sort: m.Get("sort").Int(),
			Resource: &meta.Resource{
				Host:   m.Get("host").String(),
				Path:   m.Get("path").String(),
				Method: m.Get("method").String(),
			},
			Permission: &meta.Permission{
				AuthorizedRoles: m.Get("authorized_roles").Slice().String(),
				ForbiddenRoles:  m.Get("forbidden_roles").Slice().String(),
				AllowAnyone:     m.Get("allow_anyone").Bool(),
			},
		}
		rules = append(rules, rule)
	}

	return
}
