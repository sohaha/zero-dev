package account

import (
	"zlsapp/app/error_code"
	"zlsapp/common"
	"zlsapp/service"

	"zlsapp/app/model"

	"github.com/sohaha/zlsgo/zcache"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/speps/go-hashids/v2"
	"github.com/zlsgo/zdb/builder"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	service.App
	// failedCache 防止登录爆破
	failedCache *zcache.Table
	Path        string
	Model       *model.Model
	Handlers    *AccountHandlers
}

// isBusyLogin 短时间内登录失败超过指定次数禁止登录
func (h *Account) isBusyLogin(c *znet.Context) bool {
	ip := c.GetClientIP()
	total, _ := h.failedCache.GetInt(ip)
	return total >= 5
}

func (h *Account) loginFailed(c *znet.Context) {
	ip := c.GetClientIP()
	total, _ := h.failedCache.GetInt(ip)
	data := total + 1
	h.failedCache.Set(ip, data, uint(60*data/2))
}

func (h *Account) logout(user ztype.Map) error {
	salt := zstring.Rand(8)
	_ = user.Set("salt", salt)
	return h.Handlers.Update(user.Get(model.IDKey).Value(), map[string]interface{}{
		"salt": salt,
		// "updated_at": time.Now(),
		"login_time": ztime.Time(),
	})
}

func (h *Account) Init(r *znet.Engine) {
	h.Model, _ = model.GetModel("account")

	h.failedCache = zcache.New("__account" + h.Model.Table.Name + "Failed__")

	_, _ = h.Di.Invoke(func(hashid *hashids.HashID) {
		h.Handlers = &AccountHandlers{
			Model:  h.Model,
			hashid: hashid,
		}
	})
}

// PostLogin 用户登录
func (h *Account) PostLogin(c *znet.Context) error {

	c.SetHeader("Test", "test")
	if h.isBusyLogin(c) {
		return error_code.InvalidInput.Text("错误次数过多，请稍后再试")
	}

	json, _ := c.GetJSONs()
	account := json.Get("account").String()
	password := json.Get("password").String()
	// 兼容旧版本
	if account == "" {
		account = json.Get("username").String()
	}

	if account == "" {
		return error_code.InvalidInput.Text("请输入账号")
	}

	if password == "" {
		return error_code.InvalidInput.Text("请输入密码")
	}

	user, err := h.Model.FindOne(func(b *builder.SelectBuilder) error {
		b.Where(b.EQ("account", account))
		return nil
	}, false)
	if user.IsEmpty() {
		return error_code.InvalidInput.Text("账号或密码错误")
	}
	if err != nil {
		return err
	}

	if user.IsEmpty() {
		return error_code.InvalidInput.Text("账号或密码错误")
	}
	userPassword := user.Get("password").String()

	err = bcrypt.CompareHashAndPassword(zstring.String2Bytes(userPassword), zstring.String2Bytes(password))

	if err != nil {
		h.loginFailed(c)
		return error_code.InvalidInput.Text("账号或密码错误")
	}

	status := user.Get("status").Int()
	if status != 1 {
		switch status {
		case 0:
			return error_code.Unavailable.Text("用户待激活")
		default:
			return error_code.Unavailable.Text("用户已停用")
		}
	}

	conf := ztype.Map(h.App.Conf.Core().GetStringMap("account"))

	if conf.Get("only").Bool() {
		// 	// 新登录之后清除该用户的其他端登录状态
		err = h.logout(user)
	} else {
		err = h.Handlers.Update(user.Get(model.IDKey).Value(), map[string]interface{}{
			"login_time": ztime.Time(),
		})
	}

	if err != nil {
		return err
	}

	token, _ := h.Handlers.CreateManageToken(user, conf.Get("key").String(), conf.Get("expire").Int())

	return error_code.Success.Result(c, map[string]interface{}{
		"token": token,
	})
}

// GetMessage 获取站内消息
func (h *Account) GetMessage(c *znet.Context) (interface{}, error) {
	return ztype.Map{
		"unread": 0,
	}, nil
}

// GetMe 获取当前用户信息
func (h *Account) GetMe(c *znet.Context) (interface{}, error) {
	uid := common.GetUID(c)
	info, err := h.Model.FindOne(func(b *builder.SelectBuilder) error {
		b.Where(b.EQ(model.IDKey, uid))
		b.Select(h.Model.GetFields("password", "salt")...)
		return nil
	}, false)
	return info, err
}

// // AnyLogout 用户退出
// func (h *Account) AnyLogout(c *znet.Context) error {
// 	u, ok := c.Value("user")
// 	if !ok {
// 		return restapi.ErrorMsg(restapi.ServerError, "未登录")
// 	}

// 	_ = h.logout(h.MDB, u.(ztype.Map))

// 	return restapi.Success.Result(c, nil)
// }

// // PatchPassword 修改密码
// func (h *Account) AnyPassword(c *znet.Context) error {
// 	var (
// 		old    string
// 		passwd string
// 	)
// 	rule := c.ValidRule().Required()
// 	err := zvalid.Batch(
// 		zvalid.BatchVar(&old, c.Valid(rule, "old_password", "旧密码")),
// 		zvalid.BatchVar(&passwd, c.Valid(rule, "password", "新密码")),
// 	)
// 	if err != nil {
// 		return restapi.InvalidInput.Error(err)
// 	}

// 	u, _ := c.Value("user")
// 	user := ztype.ToMap(u)
// 	if user.IsEmpty() {
// 		return restapi.InvalidInput.Error(err)
// 	}

// 	if !common.PasswordVerify(old, user.Get("password").String()) {
// 		return restapi.InvalidInput.Text("原密码错误")
// 	}

// 	key := zstring.Rand(8)
// 	_ = user.Set("key", key)
// 	newPasswd, _ := common.PasswordHash(passwd)
// 	err = h.Handlers.Update(h.MDB, user.Get("_id").Value(), map[string]interface{}{
// 		"key":      key,
// 		"password": newPasswd,
// 	})
// 	if err != nil {
// 		return restapi.ServerError.Error(err)
// 	}

// 	tokenKey := h.Conf.Auth.Key
// 	tokenExpire := h.Conf.Auth.Expire
// 	h.Handlers.ResetManageToken(c, user, tokenKey, tokenExpire)

// 	return restapi.Success.Result(c, nil)
// }
