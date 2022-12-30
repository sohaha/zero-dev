package service

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"zlsapp/internal/error_code"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zutil"
)

type (
	// App 控制器关联对象
	App struct {
		Di   zdi.Injector
		Conf *Conf
		Log  *zlog.Logger
	}
	// Router 控制器函数
	Router interface {
		Init(r *znet.Engine)
	}
	RouterAfter func(r *znet.Engine, app *App)
	Template    struct {
		DIR    string
		Global ztype.Map
	}
)

// InitWeb 初始化 WEB
func InitWeb(app *App, middlewares []znet.Handler) *znet.Engine {
	r := znet.New()
	r.Log = app.Log
	zlog.Log = r.Log

	r.BindStructSuffix = ""
	r.BindStructDelimiter = "-"
	r.SetAddr(app.Conf.Base.Port)

	isDebug := app.Conf.Base.Debug
	if isDebug {
		r.SetMode(znet.DebugMode)
	} else {
		r.SetMode(znet.ProdMode)
	}

	r.Use(znet.RewriteErrorHandler(func(c *znet.Context, err error) {

		var code int32
		statusCode := http.StatusInternalServerError
		switch zerror.GetTag(err) {
		case zerror.Internal:
			statusCode = http.StatusInternalServerError
			code = int32(error_code.ServerError)
		case zerror.InvalidInput:
			statusCode = http.StatusBadRequest
			code = int32(error_code.InvalidInput)
		case zerror.PermissionDenied:
			statusCode = http.StatusForbidden
			code = int32(error_code.PermissionDenied)
		case zerror.Unauthorized:
			statusCode = http.StatusUnauthorized
			code = int32(error_code.Unauthorized)
		default:
			errCode, ok := zerror.UnwrapCode(err)
			if ok && errCode != 0 {
				code = int32(errCode)
			} else {
				code = int32(error_code.ServerError)
			}
		}

		errMsg := strings.Join(zerror.UnwrapErrors(err), ": ")
		if errMsg == "" {
			errMsg = "unknown error"
		}

		c.JSON(int32(statusCode), map[string]interface{}{
			"code": code,
			"msg":  errMsg,
		})
	}))

	for _, middleware := range middlewares {
		r.Use(middleware)
	}

	return r
}

func RunWeb(r *znet.Engine, app *App, controllers []Router) {
	builtInRouter(r, app)

	_, err := app.Di.Invoke(func(after RouterAfter) {
		after(r, app)
	})
	if err != nil && !strings.Contains(err.Error(), "value not found for type service.RouterAfter") {
		zerror.Panic(err)
	}

	for _, c := range controllers {
		err := zutil.TryCatch(func() error {
			typeOf := reflect.TypeOf(c).Elem()
			controller := strings.TrimPrefix(typeOf.String(), "controller.")
			controller = strings.Replace(controller, ".", "/", -1)
			api := -1
			for i := 0; i < typeOf.NumField(); i++ {
				if typeOf.Field(i).Type.String() == "service.App" {
					api = i
					break
				}
			}
			if api == -1 {
				return fmt.Errorf("%s not a legitimate controller", controller)
			}

			reflect.ValueOf(c).Elem().Field(api).Set(reflect.ValueOf(*app))

			name := ""
			cName := reflect.Indirect(reflect.ValueOf(c)).FieldByName("Path")

			if cName.IsValid() && cName.String() != "" {
				name = zstring.CamelCaseToSnakeCase(cName.String(), "/")
			} else {
				name = zstring.CamelCaseToSnakeCase(controller, "/")
				if name == "home" {
					name = ""
				}
			}

			return r.BindStruct(name, c)
		})
		zerror.Panic(err)
	}

	znet.Run()
}

func StopWeb(_ *znet.Engine, _ *App) {
	znet.SetShutdown(func() {

	})
}

func builtInRouter(r *znet.Engine, app *App) {
	r.NotFoundHandler(func(c *znet.Context) {
		if c.Request.URL.Path == "/" {
			c.JSON(http.StatusOK, znet.ApiData{Code: 0, Msg: "Success", Data: ztype.Map{
				"now": ztime.Now(),
			}})
			return
		}
		c.JSON(http.StatusNotFound, znet.ApiData{Code: 404, Msg: "此路不通"})
	})
}
