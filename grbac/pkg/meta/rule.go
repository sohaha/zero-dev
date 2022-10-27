package meta

import (
	"github.com/hashicorp/go-multierror"
	jsoniter "github.com/json-iterator/go"
)

type MatchMode uint

const (
	MatchSomeAllow MatchMode = iota
	MatchAllAllow
)

// Rules is the list of Rule
type Rules []*Rule

// Rule is used to define the relationship between "resource" and "permission"
type Rule struct {
	ID          int `json:"id" yaml:"id"`
	*Resource   `yaml:",inline"`
	*Permission `yaml:",inline"`
}

// IsValid is used to test the validity of the Rule
func (rule *Rule) IsValid() error {
	if rule.Resource == nil || rule.Permission == nil {
		return ErrEmptyStructure
	}
	err := rule.Resource.IsValid()
	if err != nil {
		return err
	}
	return rule.Permission.IsValid()
}

// IsValid is used to test the validity of the Rule
func (rules Rules) IsValid() error {
	var errs error
	for _, rule := range rules {
		err := rule.IsValid()
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}

// IsRolesGranted is used to determine whether the current role is admitted by the current rule.
func (rules Rules) IsRolesGranted(roles []string, mode MatchMode) (PermissionState, error) {
	if len(rules) == 0 {
		return PermissionNeglected, nil
	}

	tail := rules[0]

	switch mode {
	case MatchAllAllow:
		for _, rule := range rules {
			state, err := rule.IsGranted(roles)
			if err != nil {
				return PermissionUngranted, err
			}
			if state != PermissionGranted {
				return PermissionUngranted, nil
			}
		}

		return PermissionGranted, nil
	default:
		for i := 0; i < len(rules); i++ {
			if tail.ID <= rules[i].ID {
				tail = rules[i]
			}
		}
		return tail.IsGranted(roles)
	}
}

func (rules Rules) String() string {
	s, _ := jsoniter.MarshalToString(rules)
	return s
}
