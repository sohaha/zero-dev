package controller

import (
	"zlsapp/common/restapi"
	"zlsapp/conf"
	"zlsapp/internal/account"
	"zlsapp/service"
)

func ManageRouter() []service.Router {
	prefix := conf.ManageRouterPrefix

	return []service.Router{
		&account.Account{
			Path: prefix + "/base",
		},
		&account.Roles{
			Path: prefix + "/account/roles",
		},
		&restapi.ManageRestApi{
			Path: prefix + "/model",
		},
	}
}
