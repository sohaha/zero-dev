package grbac

import (
	"zlsapp/grbac/meta"
)

type Resource = meta.Resource

type PermissionState = meta.PermissionState

type Permissions = meta.Permissions

type Permission = meta.Permission

type Rules = meta.Rules

type Rule = meta.Rule

type RulesLoader struct {
	rules meta.Rules
}

func NewRulesLoader(rules meta.Rules) (*RulesLoader, error) {
	return &RulesLoader{
		rules: rules,
	}, nil
}

func (loader *RulesLoader) Load() (meta.Rules, error) {
	return loader.rules, nil
}
