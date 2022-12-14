package error_code

import (
	"net/http"

	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/znet"
)

// #### 错误码为 5 位数

// | 1            | 01           | 01         |
// | :----------- | :----------- | :--------- |
// | 服务级错误码 | 模块级错误码 | 具体错误码 |

// - 服务级错误码：1 位数进行表示，比如 1 为系统错误；2 为普通错误，通常是由用户非法操作引起。
// - 模块级错误码：2 位数进行表示，比如 01 为用户模块；02 为系统模块。
// - 具体的错误码：2 位数进行表示，比如 01 为账号不存在；02 为手机号不合法。

type ErrCode zerror.ErrCode

const (
	Success ErrCode = 0

	ServerError ErrCode = 10000
	NotFound    ErrCode = 10001

	InvalidInput      ErrCode        = 20000
	UnknownClient     ErrCode        = 20001
	Unauthorized      ErrCode        = 20100
	Unauthorized2     zerror.ErrCode = 20100
	AuthorizedExpires ErrCode        = 20101
	PermissionDenied  ErrCode        = 20102
	Unavailable       ErrCode        = 20103
	InvalidAccount    ErrCode        = 20104
)

func ErrorMsg(code ErrCode, text string, err ...error) error {
	var tags []zerror.External

	switch code {
	case Unauthorized, AuthorizedExpires:
		tags = append(tags, zerror.WrapTag(zerror.Unauthorized))
	case PermissionDenied:
		tags = append(tags, zerror.WrapTag(zerror.PermissionDenied))
	case InvalidInput:
		tags = append(tags, zerror.WrapTag(zerror.InvalidInput))
	}

	if len(err) > 0 {
		return zerror.Wrap(err[0], zerror.ErrCode(code), text, tags...)
	}

	return zerror.New(zerror.ErrCode(code), text, tags...)
}

type ApiData struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg,omitempty"`
	Code ErrCode     `json:"code"`
}

func (code ErrCode) Text(msg string, err ...error) error {
	return ErrorMsg(code, msg, err...)
}

func (code ErrCode) Error(err error) error {
	return ErrorMsg(code, err.Error())
}

func (code ErrCode) Result(c *znet.Context, data interface{}, err ...error) error {
	Result(c, code, data, err...)
	return nil
}

func (n ErrCode) String() string {
	l, _ := GetI18n(n)
	return l
}

func Result(c *znet.Context, code ErrCode, data interface{}, err ...error) {
	defer c.Abort()

	var (
		m string
		d interface{} = struct{}{}
	)

	if code == Success {
		if data != nil {
			d = data
		}
		c.JSON(200, ApiData{Code: code, Msg: "", Data: d})
		return
	}

	isDebug := c.Engine.IsDebug()
	var info interface{} = struct{}{}
	if isDebug && len(err) > 0 && err[0] != nil {
		info = []string{err[0].Error()}
	}

	{
		switch v := data.(type) {
		case *zerror.Error:
			msg := zerror.UnwrapErrors(v)
			m = msg[0]
			if isDebug && len(msg) > 1 {
				switch v := info.(type) {
				case []string:
					info = append(v, msg[1:]...)
				default:
					info = msg[1:]
				}
			}
		case error:
			msg := v.Error()
			if len(msg) > 0 {
				m = msg
			}
		case string:
			if len(v) > 0 {
				m = v
			}
		}

		if len(m) == 0 {
			m, _ = GetI18n(code)
		}
	}

	var status int32 = http.StatusBadRequest
	switch true {
	case code >= 10000 && code <= 19999:
		status = http.StatusInternalServerError
	case code == 0:
		status = http.StatusOK
	default:
		switch code {
		case Unauthorized, AuthorizedExpires:
			status = http.StatusUnauthorized
		case PermissionDenied:
			status = http.StatusForbidden
		}
	}

	c.JSON(status, ApiData{Code: code, Msg: m, Data: info})
}
