package grbac

import (
	"net/http"
	"zlsapp/grbac/meta"
)

type Query struct {
	meta.Query
	rules     meta.Rules
	matchMode meta.MatchMode
}

func (q *Query) IsAllowAnyone() bool {
	return q.rules.IsAllowAnyone(q.matchMode)
}

func (q *Query) IsRolesGranted(roles []string) (PermissionState, error) {
	return q.rules.IsRolesGranted(roles, q.matchMode)
}

func (c *Engine) NewQueryByRequest(r *http.Request) (q *Query, err error) {
	if r.URL == nil {
		return nil, ErrInvalidRequest
	}

	q = &Query{
		matchMode: c.matchMode,
		Query: meta.Query{
			Path:   r.URL.Path,
			Host:   r.Host,
			Method: r.Method,
		},
	}

	q.rules, err = c.find(q)

	return q, err
}
