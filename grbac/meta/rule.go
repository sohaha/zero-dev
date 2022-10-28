package meta

import (
	"github.com/hashicorp/go-multierror"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zstring"
)

type MatchMode uint

const (
	// MatchPriorityAllow The same priority, high priority rules allowed by allowed by all
	MatchPriorityAllow MatchMode = iota
	// MatchPrioritySomeAllow The same priority, as long as there is a permissions allow, can pass
	MatchPrioritySomeAllow
	// MatchAllAllow Ignore the priority, all permissions are allowed to pass
	MatchAllAllow
)

// Rules is the list of Rule
type Rules []*Rule

// Rule is used to define the relationship between "resource" and "permission"
type Rule struct {
	Sort        int `json:"sort" yaml:"sort"`
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
	l := len(rules)
	if l == 0 {
		return PermissionNeglected, nil
	}

	tail := rules[0]

	if l == 1 {
		return tail.IsGranted(roles)
	}

	switch mode {
	case MatchAllAllow:
		for i := range rules {
			state, err := rules[i].IsGranted(roles)
			if err != nil {
				return PermissionUngranted, err
			}
			if state != PermissionGranted {
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
