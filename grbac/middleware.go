package grbac

import (
	"net/http"
	"zlsapp/app/error_code"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
)

func QueryRolesByHeaders(header http.Header) (roles []string, err error) {
	zlog.Debug("获取角色")
	roles = []string{"admin"}
	return roles, err
}

func NewMiddleware(loaderOptions ControllerOption, options ...ControllerOption) znet.Handler {
	rbac, err := New(loaderOptions, options...)
	zerror.Panic(err)

	return func(c *znet.Context) error {
		roles, err := QueryRolesByHeaders(c.Request.Header)
		if err != nil {
			return error_code.Unauthorized.Text(err.Error())
		}

		state, err := rbac.IsRequestGranted(c.Request, roles)
		if err != nil {
			return error_code.ServerError.Error(err)
		}

		if !state.IsGranted() {
			return error_code.PermissionDenied.Text("权限不足")
		}

		zlog.Debug("ok")
		c.Next()
		return nil
	}
}
