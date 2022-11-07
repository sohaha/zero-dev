package error_code

var i18n = map[string]map[ErrCode]string{
	"zh": {
		ServerError:       "内部服务器错误",
		InvalidInput:      "无效的输入",
		Unauthorized:      "未授权",
		AuthorizedExpires: "授权过期",
		PermissionDenied:  "权限不足",
		NotFound:          "资源不存在",
		Unavailable:       "不可用",
		UnknownClient:     "非法设备",
		InvalidAccount:    "无效账号",
	},
}

var DefaultLang = "zh"

func SetI18n(n map[ErrCode]string, lang ...string) {
	l := DefaultLang
	if len(lang) > 0 {
		l = lang[0]
	}
	for c, v := range n {
		if _, ok := i18n[l]; !ok {
			i18n[l] = map[ErrCode]string{}
		}
		i18n[l][c] = v
	}
}

func GetI18n(n ErrCode, lang ...string) (string, bool) {
	l := DefaultLang
	if len(lang) > 0 {
		l = lang[0]
	}
	if _, ok := i18n[l]; !ok {
		i18n[l] = map[ErrCode]string{}
	}

	t, ok := i18n[l][n]
	if !ok {
		t = "no defined"
	}
	return t, ok
}
