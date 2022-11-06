package grbac

import (
	"zlsapp/grbac/meta"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/conf"
	"gopkg.in/yaml.v3"
)

// FileLoader implements the Loader interface
// it is used to load configuration from a local file.
type FileLoader struct {
	path string
}

// FileLoader is used to initialize a FileLoader
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

// Load is used to return a list of rules
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
	zlog.Debug(rules)
	return
	// rules := meta.Rules{}
	// err = c.Unmarshal(&rules, func(dc *mapstructure.DecoderConfig) {

	// })
	// zlog.Debug(c.GetAll())
	// zlog.Debug(11, rules, err)
	// if err != nil {
	// 	return nil, err
	// }

	// return rules, nil
}

// YAMLLoader implements the Loader interface
// it is used to load configuration from a local yaml file.
type YAMLLoader struct {
	path string
}

// NewYAMLLoader is used to initialize a YAMLLoader
func NewYAMLLoader(file string) (*YAMLLoader, error) {
	loader := &YAMLLoader{
		path: zfile.RealPath(file),
	}
	_, err := loader.Load()
	if err != nil {
		return nil, err
	}
	return loader, nil
}

// Load is used to return a list of rules
func (loader *YAMLLoader) Load() (meta.Rules, error) {
	bytes, err := zfile.ReadFile(loader.path)
	if err != nil {
		return nil, err
	}
	rules := meta.Rules{}
	err = yaml.Unmarshal(bytes, &rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// JSONLoader implements the Loader interface
// it is used to load configuration from a local json file.
type JSONLoader struct {
	path string
}

// NewJSONLoader is used to initialize a JSONLoader
func NewJSONLoader(file string) (*JSONLoader, error) {
	loader := &JSONLoader{
		path: zfile.RealPath(file),
	}
	_, err := loader.Load()
	if err != nil {
		return nil, err
	}
	return loader, nil
}

// Load is used to return a list of rules
func (loader *JSONLoader) Load() (meta.Rules, error) {
	bytes, err := zfile.ReadFile(loader.path)
	if err != nil {
		return nil, err
	}
	rules := meta.Rules{}
	err = zjson.Unmarshal(bytes, &rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}
