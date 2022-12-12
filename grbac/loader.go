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

	rules = ParseMap(c.GetAll())

	return
}

func ParseMap(m map[string]interface{}) meta.Rules {
	rules := make(meta.Rules, 0)
	for _, v := range m {
		m := ztype.ToMap(v)
		host := m.Get("host").String()
		if host == "" {
			host = "*"
		}
		method := m.Get("method").String()
		if method == "" {
			method = "*"
		}
		rule := &meta.Rule{
			Sort: m.Get("sort").Int(),
			Resource: &meta.Resource{
				Host:   host,
				Path:   m.Get("path").String(),
				Method: method,
			},
			Permission: &meta.Permission{
				AuthorizedRoles: m.Get("authorized_roles").Slice().String(),
				ForbiddenRoles:  m.Get("forbidden_roles").Slice().String(),
				AllowAnyone:     m.Get("allow_anyone").Bool(),
			},
		}
		rules = append(rules, rule)
	}
	return rules
}
