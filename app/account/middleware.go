package account

import (
	"time"
	"zlsapp/app/error_code"
	"zlsapp/grbac"
	"zlsapp/grbac/meta"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/speps/go-hashids/v2"
)

func NewMiddleware(app *service.App) znet.Handler {
	loaderOptions := grbac.WithFile("grbac/testdata/grbac", time.Second*2)
	options := grbac.WithMatchMode(meta.MatchPrioritySomeAllow)
	rbac, err := grbac.New(loaderOptions, options)
	zerror.Panic(err)

	pubPath := []string{"/base/login"}

	m, err := migration(app.Di)
	zerror.Panic(err)

	zlog.Debug(m, err)
	h := &AccountHandlers{
		Model: m,
	}
	_, _ = app.Di.Invoke(func(hashid *hashids.HashID) {
		h.hashid = hashid
	})

	key := app.Conf.Core().GetString("account.key")

	return func(c *znet.Context) error {
		path := c.Request.URL.Path
		if zarray.Contains(pubPath, path) {
			c.Next()
			return nil
		}

		j, err := h.ParsingManageToken(c, key)
		if err != nil {
			return err
		}

		uid, roles, err := h.QueryRoles(j)
		if err != nil {
			return error_code.Unauthorized.Text(err.Error())
		}
		c.WithValue("uid", uid)

		state, err := rbac.IsRequestGranted(c.Request, roles)
		if err != nil {
			return error_code.ServerError.Error(err)
		}

		if !state.IsGranted() {
			return error_code.PermissionDenied.Text("权限不足")
		}

		c.Next()
		return nil
	}
}
