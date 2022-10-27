package grbac

import (
	"zlsapp/grbac/pkg/meta"
)

// AdvancedRule allows you to write RBAC rules in a more concise way
type AdvancedRule struct {
	Host   []string `json:"host"`
	Path   []string `json:"path"`
	Method []string `json:"method"`

	*meta.Permission
}

// AdvancedRules is the list of AdvancedRules
type AdvancedRules []*AdvancedRule

// GetRules is used to convert AdvancedRules to meta.Rules
func (adv AdvancedRules) GetRules() meta.Rules {
	var rules meta.Rules
	for _, item := range adv {
		for _, host := range item.Host {
			for _, path := range item.Path {
				for _, method := range item.Method {
					rules = append(rules, &meta.Rule{
						Resource: &meta.Resource{
							Host:   host,
							Path:   path,
							Method: method,
						},
						Permission: item.Permission,
					})
				}
			}
		}
	}
	return rules
}

// AdvancedRulesLoader implements the Loader interface
// it is used to load configuration from advanced data.
type AdvancedRulesLoader struct {
	rules AdvancedRules
}

// NewAdvancedRulesLoader is used to initialize a AdvancedRulesLoader
func NewAdvancedRulesLoader(rules AdvancedRules) (*AdvancedRulesLoader, error) {
	return &AdvancedRulesLoader{
		rules: rules,
	}, nil
}

// Load is used to return a list of rules
func (loader *AdvancedRulesLoader) Load() (meta.Rules, error) {
	return loader.rules.GetRules(), nil
}
