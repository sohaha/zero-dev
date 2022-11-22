package account

import (
	"time"

	"zlsapp/grbac"
	"zlsapp/grbac/meta"
	"zlsapp/internal/error_code"
	"zlsapp/internal/parse"
	"zlsapp/service"

	_ "embed"

	"github.com/pelletier/go-toml/v2"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/speps/go-hashids/v2"
)

//go:embed permission.toml
var defPermission []byte

func NewMiddleware(app *service.App, pubPath []string) znet.Handler {
	var loaderOptions grbac.Option
	if zfile.FileExist("permission.toml") {
		loaderOptions = grbac.WithFile("permission.toml", time.Second*2)
	} else {
		loaderOptions = grbac.WithLoader(func() (grbac.Rules, error) {
			var m map[string]interface{}
			err := toml.Unmarshal([]byte(defPermission), &m)
			if err != nil {
				return nil, err
			}
			return grbac.ParseMap(m), nil
		}, -1)
	}

	options := grbac.WithMatchMode(meta.MatchPrioritySomeAllow)
	rbac, err := grbac.New(loaderOptions, options)
	zerror.Panic(err)

	m, err := migration(app.Di)
	zerror.Panic(err)

	h := &AccountHandlers{
		Model: m,
	}
	_, _ = app.Di.Invoke(func(hashid *hashids.HashID) {
		h.hashid = hashid
	})

	key := app.Conf.Core().GetString("account.key")

	return func(c *znet.Context) error {
		path := c.Request.URL.Path
		for _, v := range pubPath {
			if zstring.Match(path, v) {
				c.Next()
				return nil
			}
		}

		q, err := rbac.NewQueryByRequest(c.Request)
		if err != nil {
			return err
		}

		if q.IsAllowAnyone() {
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

		state, err := q.IsRolesGranted(roles)
		if err != nil {
			return error_code.ServerError.Error(err)
		}

		if !state.IsGranted() {
			return error_code.PermissionDenied.Text("权限不足")
		}

		c.WithValue("uid", user.Get(parse.IDKey).Value())
		c.Next()
		return nil
	}
}
