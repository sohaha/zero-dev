package grbac

import (
	"io/ioutil"

	"zlsapp/grbac/meta"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zjson"
	"gopkg.in/yaml.v3"
)

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
	bytes, err := ioutil.ReadFile(loader.path)
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
	bytes, err := ioutil.ReadFile(loader.path)
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
