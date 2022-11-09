package grbac

import (
	"zlsapp/grbac/meta"
)

type Result struct {
	Error error
	State PermissionState
}

func NewQuery(c *Controller, host, path, method string, roles []string) *Result {
	state, err := c.IsQueryGranted(&meta.Query{
		Host:   host,
		Path:   path,
		Method: method,
	}, roles)
	return &Result{
		State: state,
		Error: err,
	}
}
