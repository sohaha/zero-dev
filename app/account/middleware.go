package account

import (
	"time"
	"zlsapp/app/error_code"
	"zlsapp/app/model"
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

		user, err := h.QueryRoles(j)
		if err != nil {
			return error_code.Unauthorized.Text(err.Error())
		}

		if app.Conf.Core().GetBool("account.only") {
			salt := user.Get("salt").String()
			if salt != j.U[:8] {
				return error_code.AuthorizedExpires.Text("登录状态失效，请重新登录")
			}
		}

		roles := user.Get("roles").Slice().String()
		c.WithValue("uid", user.Get(model.IDKey).Value())

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
