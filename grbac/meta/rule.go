package meta

import (
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zstring"
)

type MatchMode uint

const (
	// MatchPrioritySomeAllow The same priority, as long as there is a permissions allow, can pass
	MatchPrioritySomeAllow MatchMode = iota
	// MatchPriorityAllow The same priority, high priority rules allowed by allowed by all
	MatchPriorityAllow
	// MatchAllAllow Ignore the priority, all permissions are allowed to pass
	MatchAllAllow
	// MatchSomeAllow As long as there is a permission, you can pass
	MatchSomeAllow
)

// Rules is the list of Rule
type Rules []*Rule

// Rule is used to define the relationship between "resource" and "permission"
type Rule struct {
	*Resource   `mapstructure:",squash" yaml:",inline"`
	*Permission `yaml:",inline"`
	Sort        int `mapstructure:"sort" json:"sort" yaml:"sort"`
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
			errs = zerror.With(errs, err.Error())
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}

// IsRolesGranted is used to determine whether the current role is admitted by the current rule.
func (rules Rules) IsRolesGranted(roles []string, mode MatchMode) (PermissionState, error) {
	l := len(rules)
	if l == 0 {
		return PermissionNeglected, nil
	}

	tail := rules[0]

	if l == 1 {
		return tail.IsGranted(roles)
	}

	switch mode {
	case MatchAllAllow, MatchSomeAllow:
		for i := range rules {
			state, err := rules[i].IsGranted(roles)
			if err != nil {
				return PermissionUngranted, err
			}
			ok := state == PermissionGranted
			if ok && MatchSomeAllow == mode {
				return PermissionGranted, nil
			}
			if !ok {
				return PermissionUngranted, nil
			}
		}

		return PermissionGranted, nil
	default:
		priorityRules := Rules{tail}
		for i := 1; i < l; i++ {
			if tail.Sort < rules[i].Sort {
				tail = rules[i]
				priorityRules = Rules{tail}
			} else if tail.Sort == rules[i].Sort {
				priorityRules = append(priorityRules, rules[i])
			}
		}

		for i := range priorityRules {
			state, err := priorityRules[i].IsGranted(roles)
			if err != nil {
				return PermissionUngranted, err
			}
			ok := state == PermissionGranted
			if ok && MatchPrioritySomeAllow == mode {
				return PermissionGranted, nil
			}

			if !ok {
				return PermissionUngranted, nil
			}
		}

		return PermissionGranted, nil
	}
}

func (rules Rules) String() string {
	s, _ := zjson.Marshal(rules)
	return zstring.Bytes2String(s)
}
