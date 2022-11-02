package account

import (
	"time"
	"zlsapp/app/error_code"
	"zlsapp/app/model"
	"zlsapp/common/jwt"
	"zlsapp/service"

	"github.com/sohaha/zlsgo/zcache"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
)

type Account struct {
	service.App
	// failedCache 防止登录爆破
	failedCache *zcache.Table
	Path        string
	Model       string
	Handlers    *AccountHandlers
}
type AccountHandlers struct {
	Model string
}

func (h *AccountHandlers) Update(db *zdb.DB, filter, update interface{}) error {
	// col := db.Collection(context.TODO(), h.Model)
	// row, err := col.FindOneAndUpdate(filter, update)
	// if err == nil {
	// 	_, _ = zcache.New("__account_" + h.Model + "__").Delete(zmgo.MustObjectID(row.Get("_id").Value()).Hex())
	// }
	// return err
	return nil
}

func (h *AccountHandlers) CacheForID(db *zdb.DB, uid string) (row ztype.Map, err error) {
	// data, err := zcache.New("__account_"+h.Model+"__").MustGet(uid, func(set func(data interface{},
	// 	lifeSpan time.Duration, interval ...bool)) (err error) {
	// 	row, err := db.Collection(context.TODO(), h.Model).FindOne(uid, func(opts *options.FindOneOptions) {
	// 	})
	// 	if row.IsEmpty() {
	// 		return errors.New("账号不存在:" + uid)
	// 	}

	// 	if err != nil {
	// 		return err
	// 	}

	// 	set(row, time.Second*60)
	// 	return nil
	// })
	// if err != nil {
	// 	return make(ztype.Map), err
	// }

	// m, ok := data.(ztype.Map)
	// if !ok {
	// 	return make(ztype.Map), errors.New("用户信息获取失败")
	// }

	return
}

func (h *AccountHandlers) CreateManageToken(user ztype.Map, key string, expire int) (string, error) {
	// u := user.Get("key").String() + zmgo.MustObjectID(user.Get("_id").Value()).Hex()
	// claims := jwt.JwtInfo{
	// 	U: u,
	// 	StandardClaims: gjwt.StandardClaims{
	// 		ExpiresAt: time.Now().Add(time.Duration(expire) * time.Second).Unix(),
	// 	},
	// }

	// token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, claims)
	// signedToken, err := token.SignedString(zstring.String2Bytes(key))
	// if err != nil {
	// 	return "", fmt.Errorf("生成签名失败: %v", err)
	// }
	// return signedToken, nil
	return "", nil
}

func (h *AccountHandlers) ParsingManageToken(c *znet.Context, key string) (*jwt.JwtInfo, error) {
	// token := jwt.GetToken(c)
	// if token == "" {
	// 	return nil, restapi.Unauthorized.Text("请先登录")
	// }

	// j, err := jwt.ParsingToken(token, key)
	// if err != nil {
	// 	return nil, restapi.AuthorizedExpires.Text("登录状态过期，请重新登录")
	// }

	// if len(j.U) < 12 {
	// 	return nil, restapi.InvalidInput.Text("无效签名")
	// }

	// return j, nil
	return nil, nil
}

func (h *AccountHandlers) ResetManageToken(c *znet.Context, user ztype.Map, key string, expire int) {
	token, err := h.CreateManageToken(user, key, expire)
	if err == nil {
		c.SetHeader("Re-Token", token)
	}
}

func (h *AccountHandlers) Middleware(db *zdb.DB, key string, pubRouter []string) func(c *znet.Context) (err error) {
	return func(c *znet.Context) (err error) {
		// url := c.Request.URL.Path
		// for _, v := range pubRouter {
		// 	if url == v {
		// 		c.Next()
		// 		return
		// 	}
		// }

		// j, err := h.ParsingManageToken(c, key)
		// if err != nil {
		// 	if errCode, ok := zerror.UnwrapCode(err); ok {
		// 		return restapi.ErrCode(errCode).Error(err)
		// 	}

		// 	return restapi.InvalidInput.Error(err)
		// }

		// salt := j.U[0:8]
		// uid := j.U[8:]
		// row, err := h.CacheForID(db, uid)
		// if err != nil {
		// 	return restapi.InvalidInput.Error(err)
		// }

		// userSalt := row.Get("key").String()
		// if userSalt == "" || userSalt != salt {
		// 	return restapi.AuthorizedExpires.Text("登录状态已失效")
		// }

		// c.WithValue("uid", uid)
		// c.WithValue("user", row)
		// c.Next()

		return nil
	}
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

func (h *Account) logout(db *zdb.DB, user ztype.Map) error {
	key := zstring.Rand(8)
	_ = user.Set("key", key)
	return h.Handlers.Update(db, user.Get("_id").Value(), map[string]interface{}{
		"key": key,
		// "updated_at": time.Now(),
		"login_time": time.Now(),
	})
}

func (h *Account) Init(r *znet.Engine) {
	h.failedCache = zcache.New("__account" + h.Model + "Failed__")
	h.Handlers = &AccountHandlers{Model: h.Model}

	_, err := h.Di.Invoke(func(db *zdb.DB) {
		json, _ := zjson.SetBytes([]byte("{}"), "name", ztype.Map{})
		json, _ = zjson.SetBytes(json, "table", ztype.Map{
			"name":    "account_user",
			"comment": "用户表",
		})

		json, _ = zjson.SetBytes(json, "options", ztype.Map{
			"timestamps": true,
		})

		json, _ = zjson.SetBytes(json, "values", ztype.Maps{
			{
				"account":  "admin",
				"password": "admin",
			},
		})
		zlog.Success("初始化用户")
		zlog.Printf("        账号: %s\n        密码: %s\n", "admin", "admin")
		json, _ = zjson.SetBytes(json, "columns", ztype.Maps{
			{
				"label":    "头像",
				"name":     "avatar",
				"nullable": true,
				"type":     "string",
				"validations": ztype.Maps{
					{
						"method": "url",
					},
				},
			},
			{
				"name":  "account",
				"type":  "string",
				"label": "账号",
				"validations": ztype.Maps{
					{
						"method": "minLength",
						"args":   3,
					},
					{
						"method": "maxLength",
						"args":   10,
					},
				},
			},
			{
				"name":  "password",
				"type":  "string",
				"label": "密码",
				"validations": ztype.Maps{
					{
						"method": "minLength",
						"args":   3,
					},
					{
						"method": "maxLength",
						"args":   20,
					},
				},
			},
		})

		m, err := model.Add(db, "account", json)
		zerror.Panic(err)
		zerror.Panic(m.Migration().Auto())
	})
	zerror.Panic(err)
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

	// col := h.MDB.Collection(context.TODO(), h.Model)

	// user, _ := col.FindOne(map[string]interface{}{"username": data.Username})

	// if user.IsEmpty() || !common.PasswordVerify(data.Password, user.Get("password").String()) {
	// 	h.loginFailed(c)
	// 	return error_code.InvalidAccount.Text("账号/密码错误")
	// }

	// status := user.Get("status").Int()
	// if status != 1 {
	// 	switch status {
	// 	case 0:
	// 		return error_code.Unavailable.Text("用户待激活")
	// 	default:
	// 		return error_code.Unavailable.Text("用户已停用")
	// 	}

	// }

	// if h.App.Conf.Auth.Only {
	// 	// 新登录之后清除该用户的其他端登录状态
	// 	_ = h.logout(h.MDB, user)
	// } else {
	// 	_ = h.Handlers.Update(h.MDB, user.Get("_id").Value(), map[string]interface{}{
	// 		// "updated_at": time.Now(),
	// 		"login_time": time.Now(),
	// 	})
	// }

	// manage := h.App.Conf.Auth
	// token, _ := h.Handlers.CreateManageToken(user, manage.Key, manage.Expire)

	// return restapi.Success.Result(c, map[string]interface{}{
	// 	"token": token,
	// 	"id":    zmgo.MustObjectID(user.Get("_id").Value()).Hex(),
	// })
	return nil
}

// GetMessage 获取站内消息
func (h *Account) GetMessage(c *znet.Context) error {
	return error_code.Success.Result(c, map[string]interface{}{
		"unread": 0,
	})
}

// // GetMe 获取当前用户信息
// func (h *Account) GetMe(c *znet.Context) error {
// 	uid := logic.GetUID(c)
// 	info, _ := h.MDB.Collection(context.TODO(), h.Model).FindOne(uid, func(opts *options.FindOneOptions) {
// 		opts.Projection = bson.M{"password": 0, "key": 0}
// 	})
// 	res := map[string]interface{}{
// 		"info": info,
// 	}
// 	roles, ok := c.Value("roles")
// 	if ok {
// 		res["roles"] = roles
// 	}

// 	return restapi.Success.Result(c, res)
// }

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
