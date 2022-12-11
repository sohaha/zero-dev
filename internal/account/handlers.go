package account

import (
	"errors"
	"fmt"
	"time"

	"zlsapp/common/hashid"
	"zlsapp/common/jwt"
	"zlsapp/internal/error_code"
	"zlsapp/internal/parse"

	gjwt "github.com/golang-jwt/jwt"
	"github.com/sohaha/zlsgo/zcache"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/speps/go-hashids/v2"
)

type AccountHandlers struct {
	Model  *parse.Modeler
	hashid *hashids.HashID
}

func (h *AccountHandlers) Update(id interface{}, update ztype.Map) error {
	filter := ztype.Map{
		parse.IDKey: id,
	}
	row, err := parse.FindOne(h.Model, filter)
	if row.IsEmpty() || err != nil {
		return errors.New("账号不存在")
	}

	_, err = parse.Update(h.Model, filter, update)
	if err != nil {
		return err
	}

	_, _ = h.Cache().Delete(row.Get(parse.IDKey).String())
	return nil
}

func (h *AccountHandlers) Cache() *zcache.Table {
	return zcache.New("__account_" + h.Model.Table.Name + "__")
}

// CacheForID 从缓存中获取用户信息
func (h *AccountHandlers) CacheForID(uid interface{}) (row ztype.Map, err error) {
	idStr := ztype.ToString(uid)
	data, err := h.Cache().MustGet(idStr, func(set func(data interface{},
		lifeSpan time.Duration, interval ...bool)) (err error) {
		row, err := parse.FindOne(h.Model, ztype.Map{parse.IDKey: uid})

		if row.IsEmpty() {
			return errors.New("账号不存在: " + idStr)
		}

		if err != nil {
			return err
		}

		set(row, time.Second*60)
		return nil
	})
	if err != nil {
		return make(ztype.Map), err
	}

	m, ok := data.(ztype.Map)
	if !ok {
		return make(ztype.Map), errors.New("用户信息获取失败")
	}

	return m, nil
}

func (h *AccountHandlers) CreateManageToken(user ztype.Map, key string, expire int) (string, error) {
	var id string
	var err error

	if h.Model.Options.CryptID {
		i := user.Get(parse.IDKey).Int64()
		if i == 0 {
			id = user.Get(parse.IDKey).String()
		} else {
			id, err = hashid.EncryptID(h.hashid, i)
			if err != nil {
				return "", err
			}
		}
	} else {
		id = user.Get(parse.IDKey).String()
	}

	if id == "" {
		return "", errors.New("用户ID不能为空")
	}

	u := user.Get("salt").String() + id

	if expire == 0 {
		expire = 3600 * 24 * 30
	}
	claims := jwt.JwtInfo{
		U: u,
		StandardClaims: gjwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expire) * time.Second).Unix(),
		},
	}

	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(zstring.String2Bytes(key))
	if err != nil {
		return "", fmt.Errorf("生成签名失败: %v", err)
	}

	return signedToken, nil
}

func (h *AccountHandlers) ParsingManageToken(c *znet.Context, key string) (*jwt.JwtInfo, error) {
	token := jwt.GetToken(c)
	if token == "" {
		return nil, zerror.New(error_code.Unauthorized2, "请先登录", zerror.WrapTag(zerror.Unauthorized))
	}

	j, err := jwt.ParsingToken(token, key)
	if err != nil {
		return nil, error_code.AuthorizedExpires.Text("登录状态过期，请重新登录")
	}

	if len(j.U) < 8 {
		return nil, error_code.InvalidInput.Text("无效签名")
	}

	return j, nil
}

func (h *AccountHandlers) ResetManageToken(c *znet.Context, user ztype.Map, key string, expire int) string {
	token, err := h.CreateManageToken(user, key, expire)
	if err == nil {
		c.SetHeader("Re-Token", token)
	}
	return token
}

func (h *AccountHandlers) QueryRoles(j *jwt.JwtInfo) (user ztype.Map, err error) {
	var uid interface{}
	if h.Model.Options.CryptID {
		var id int64
		id, err = hashid.DecryptID(h.hashid, j.U[8:])
		if err != nil {
			return nil, errors.New("无效签名")
		}
		uid = ztype.ToString(id)
	} else {
		uid = j.U[8:]
	}

	user, err = h.CacheForID(uid)

	return
}
