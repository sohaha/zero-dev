package meta

import (
	"github.com/sohaha/zlsgo/zerror"
)

type Permissions []*Permission

type Permission struct {
	AuthorizedRoles []string
	ForbiddenRoles  []string
	AllowAnyone     bool
}

func (p *Permission) IsValid() error {
	if !p.AllowAnyone && len(p.AuthorizedRoles) == 0 && len(p.ForbiddenRoles) == 0 {
		return zerror.With(ErrEmptyStructure, "permission: ")
	}
	return nil
}

func (p *Permission) IsGranted(roles []string) (PermissionState, error) {
	if p.AllowAnyone {
		return PermissionGranted, nil
	}

	if len(roles) == 0 {
		return PermissionUngranted, nil
	}

	for _, role := range roles {
		for _, forbidden := range p.ForbiddenRoles {
			if forbidden == "*" || (role == forbidden) {
				return PermissionUngranted, nil
			}
		}
		for _, authorized := range p.AuthorizedRoles {
			if authorized == "*" || (role == authorized) {
				return PermissionGranted, nil
			}
		}
	}
	return PermissionUngranted, nil
}
