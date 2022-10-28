package grbac

import (
	"zlsapp/grbac/meta"
)

// Resource defines resources
type Resource = meta.Resource

// PermissionState identifies the status of the permission
type PermissionState = meta.PermissionState

// Permissions is the set of Permission
type Permissions = meta.Permissions

// Permission is used to define permission control information
type Permission = meta.Permission

// Rules is the list of Rule
type Rules = meta.Rules

// Rule is used to define the relationship between "resource" and "permission"
type Rule = meta.Rule

// Query defines the data structure of the query parameters
type Query = meta.Query

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
