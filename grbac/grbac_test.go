package grbac

import (
	"zlsapp/grbac/meta"
)

type Result struct {
	State PermissionState
	Error error
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
