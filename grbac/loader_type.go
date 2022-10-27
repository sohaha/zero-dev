package grbac

import (
	"zlsapp/grbac/pkg/meta"
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
