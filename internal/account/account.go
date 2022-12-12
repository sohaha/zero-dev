package account

import (
	"errors"
	"zlsapp/common"
	"zlsapp/internal/error_code"
	"zlsapp/service"

	"zlsapp/internal/parse"

	"github.com/sohaha/zlsgo/zcache"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/sohaha/zlsgo/zvalid"
	"github.com/speps/go-hashids/v2"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	service.App
	failedCache *zcache.Table
	Model       *parse.Modeler
	Handlers    *AccountHandlers
	Path        string
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
	return h.Handlers.Update(user.Get(parse.IDKey).Value(), map[string]interface{}{
		"salt": salt,
		// "updated_at": time.Now(),
		"login_time": ztime.Time(),
	})
}

func (h *Account) Init(r *znet.Engine) {
	var ok bool
	h.Model, ok = parse.GetModel(UsersModel)
	if !ok {
		zerror.Panic(errors.New("model account not found"))
	}

	h.failedCache = zcache.New("__account" + h.Model.Table.Name + "Failed__")

	_, _ = h.Di.Invoke(func(hashid *hashids.HashID) {
		h.Handlers = &AccountHandlers{
			Model:  h.Model,
			hashid: hashid,
		}
	})
}

// PostLogin 用户登录
func (h *Account) PostLogin(c *znet.Context) (any, error) {
	c.SetHeader("Test", "test")
	if h.isBusyLogin(c) {
		return nil, error_code.InvalidInput.Text("错误次数过多，请稍后再试")
	}

	json, _ := c.GetJSONs()
	account := json.Get("account").String()
	password := json.Get("password").String()
	// 兼容旧版本
	if account == "" {
		account = json.Get("username").String()
	}

	if account == "" {
		return nil, error_code.InvalidInput.Text("请输入账号")
	}

	if password == "" {
		return nil, error_code.InvalidInput.Text("请输入密码")
	}

	user, err := parse.FindOne(h.Model, ztype.Map{
		"account": account,
	})
	if user.IsEmpty() {
		return nil, error_code.InvalidInput.Text("账号或密码错误")
	}
	if err != nil {
		return nil, err
	}

	userPassword := user.Get("password").String()

	err = bcrypt.CompareHashAndPassword(zstring.String2Bytes(userPassword), zstring.String2Bytes(password))

	if err != nil {
		h.loginFailed(c)

		WrapLogs(c, "登录失败", func(data *ztype.Map) {
			(*data)["category"] = LogTypeLogin
			(*data)["detail"] = "密码错误"
			(*data)["status"] = LogsStatusUnread
			(*data)["uid"] = user.Get(parse.IDKey).String()
		})

		return nil, error_code.InvalidInput.Text("账号或密码错误")
	}

	status := user.Get("status").Int()
	if status != 1 {
		switch status {
		case 0:
			return nil, error_code.Unavailable.Text("用户待激活")
		default:
			return nil, error_code.Unavailable.Text("用户已停用")
		}
	}

	conf := ztype.Map(h.App.Conf.Core().GetStringMap("account"))

	if conf.Get("only").Bool() {
		// 	// 新登录之后清除该用户的其他端登录状态
		err = h.logout(user)
	} else {
		err = h.Handlers.Update(user.Get(parse.IDKey).Value(), map[string]interface{}{
			"login_time": ztime.Time(),
		})
	}

	if err != nil {
		return nil, err
	}

	token, _ := h.Handlers.CreateManageToken(user, conf.Get("key").String(), conf.Get("expire").Int())

	WrapLogs(c, "登录成功", func(data *ztype.Map) {
		(*data)["category"] = LogTypeLogin
		(*data)["uid"] = user.Get(parse.IDKey).String()
	})

	return map[string]interface{}{
		"token": token,
	}, nil
}

// GetMessage 获取站内消息
func (h *Account) GetMessage(c *znet.Context) (any, error) {
	m, _ := parse.GetModel(LogsModel)
	uid := common.GetUID(c)
	item, _ := parse.FindOne(m, ztype.Map{
		"uid":    uid,
		"status": LogsStatusUnread,
	}, func(so *parse.StorageOptions) error {
		so.Fields = []string{"count(*) as unread"}
		return nil
	})
	return ztype.Map{
		"unread": item.Get("unread").Int(),
	}, nil
}

// GetLogs 获取操作日志
func (h *Account) GetLogs(c *znet.Context) (any, error) {
	m, _ := parse.GetModel(LogsModel)
	// page, size, err := parse.GetPages(c)
	// if err != nil {
	// 	return nil, zerror.InvalidInput.Wrap(err, "Invalid page or size")
	// }

	filter := ztype.Map{
		"uid": common.GetUID(c),
	}

	category, _ := c.GetQuery("category")
	if category != "" {
		filter["category"] = category
	}

	// items, p, err := parse.Pages(m, page, size, filter, func(so *parse.StorageOptions) error {
	// 	so.OrderBy = map[string]int8{parse.IDKey: -1}
	// 	return nil
	// })
	fields := m.GetFields()
	withFilds := []string{"user"}

	res, err := parse.RestapiGetPage(c, m, filter, fields, withFilds, func(so *parse.StorageOptions) error {
		// so.Wheres = []parse.StorageWhere{
		// 	{
		// 		Field: "category",
		// 		Expr:  "eq",
		// 		Value: category,
		// 	},
		// }
		return nil
	})
	if err != nil {
		return nil, err
	}

	go func() {
		if res.Page.Total == 0 {
			return
		}
		ids := []interface{}{}
		for _, item := range res.Items {
			if item.Get("status").Int() == LogsStatusUnread {
				ids = append(ids, item.Get(parse.IDKey).String())
			}
		}

		if len(ids) == 0 {
			return
		}
		_, _ = parse.Update(m, ztype.Map{
			parse.IDKey: ids,
		}, ztype.Map{
			"status": LogsStatusRead,
		})
	}()

	return res, err
}

// GetMe 获取当前用户信息
func (h *Account) GetMe(c *znet.Context) (any, error) {
	info, err := parse.FindOne(h.Model, ztype.Map{
		parse.IDKey: common.GetUID(c),
	}, func(so *parse.StorageOptions) error {
		so.Fields = h.Model.GetFields("password", "salt")
		return nil
	})
	return ztype.Map{
		"info": info,
	}, err
}

// PatchMe 修改当前用户信息
func (h *Account) PatchMe(c *znet.Context) (any, error) {
	uid := common.GetUID(c)
	data, _ := c.GetJSONs()
	err := h.Handlers.Update(uid, data.MapString())
	return nil, err
}

// AnyLogout 用户退出
func (h *Account) AnyLogout(c *znet.Context) (any, error) {
	uid := common.GetUID(c)
	user, err := h.Handlers.CacheForID(uid)
	if err != nil {
		return nil, err
	}

	err = h.logout(user)
	return nil, err
}

// PatchPassword 修改密码
func (h *Account) PatchPassword(c *znet.Context) (any, error) {
	var (
		old    string
		passwd string
	)
	rule := c.ValidRule().Required()
	err := zvalid.Batch(
		zvalid.BatchVar(&old, c.Valid(rule, "old_password", "旧密码")),
		zvalid.BatchVar(&passwd, c.Valid(rule, "password", "新密码")),
	)
	if err != nil {
		return nil, error_code.InvalidInput.Error(err)
	}

	uid := common.GetUID(c)
	user, _ := h.Handlers.CacheForID(uid)
	if user.IsEmpty() {
		return nil, error_code.InvalidInput.Error(err)
	}

	if !zvalid.Text(old).CheckPassword(user.Get("password").String()).Ok() {
		return nil, error_code.InvalidInput.Text("原密码错误")
	}

	salt := zstring.Rand(8)
	_ = user.Set("salt", salt)

	err = h.Handlers.Update(uid, map[string]interface{}{
		"salt":     salt,
		"password": passwd,
	})
	if err != nil {
		return nil, error_code.ServerError.Error(err)
	}

	conf := ztype.Map(h.App.Conf.Core().GetStringMap("account"))
	token := h.Handlers.ResetManageToken(c, user, conf.Get("key").String(), conf.Get("expire").Int())

	WrapLogs(c, "修改密码")

	return ztype.Map{
		"token": token,
	}, nil
}
