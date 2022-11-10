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

type Rules []*Rule

type Rule struct {
	*Resource   `mapstructure:",squash" yaml:",inline"`
	*Permission `yaml:",inline"`
	Sort        int `mapstructure:"sort" json:"sort" yaml:"sort"`
}

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

func (rules Rules) IsAllowAnyone(mode MatchMode) bool {
	l := len(rules)
	if l == 0 {
		return false
	}

	r := rules[0]

	if l == 1 {
		return r.AllowAnyone
	}

	switch mode {
	case MatchAllAllow, MatchSomeAllow:
		for i := range rules {
			ok := rules[i].AllowAnyone
			if ok && MatchSomeAllow == mode {
				return true
			}
			if !ok {
				return false
			}
		}

		return true
	default:
		priorityRules := Rules{r}
		for i := 1; i < l; i++ {
			if r.Sort < rules[i].Sort {
				r = rules[i]
				priorityRules = Rules{r}
			} else if r.Sort == rules[i].Sort {
				priorityRules = append(priorityRules, rules[i])
			}
		}

		for i := range priorityRules {
			ok := priorityRules[i].AllowAnyone
			if ok && MatchPrioritySomeAllow == mode {
				return true
			}

			if !ok {
				return false
			}
		}

		return true
	}
}

func (rules Rules) IsRolesGranted(roles []string, mode MatchMode) (PermissionState, error) {
	l := len(rules)
	if l == 0 {
		return PermissionNeglected, nil
	}

	r := rules[0]

	if l == 1 {
		return r.IsGranted(roles)
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
		priorityRules := Rules{r}
		for i := 1; i < l; i++ {
			if r.Sort < rules[i].Sort {
				r = rules[i]
				priorityRules = Rules{r}
			} else if r.Sort == rules[i].Sort {
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
