package grbac

import (
	"zlsapp/grbac/pkg/meta"
)

// RulesLoader implements the Loader interface
// it is used to load configuration from given rules.
type RulesLoader struct {
	rules meta.Rules
}

// NewRulesLoader is used to initialize a RulesLoader
func NewRulesLoader(rules meta.Rules) (*RulesLoader, error) {
	return &RulesLoader{
		rules: rules,
	}, nil
}

// Load is used to return a list of rules
func (loader *RulesLoader) Load() (meta.Rules, error) {
	return loader.rules, nil
}
