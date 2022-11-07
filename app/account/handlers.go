package account

import (
	"errors"
	"fmt"
	"time"
	"zlsapp/app/error_code"
	"zlsapp/app/model"
	"zlsapp/common/hashid"
	"zlsapp/common/jwt"

	gjwt "github.com/golang-jwt/jwt"
	"github.com/sohaha/zlsgo/zcache"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/speps/go-hashids/v2"
	"github.com/zlsgo/zdb/builder"
)

type AccountHandlers struct {
	Model  *model.Model
	hashid *hashids.HashID
}

func (h *AccountHandlers) Update(id interface{}, update interface{}) error {

	row, _ := h.Model.FindOne(func(b *builder.SelectBuilder) error {
		b.Where(b.EQ(model.IDKey, id))
		return nil
	}, true)

	if row.IsEmpty() {
		return errors.New("账号不存在")
	}

	_, err := h.Model.Update(update, func(b *builder.UpdateBuilder) error {
		b.Where(b.EQ(model.IDKey, id))
		return nil
	})
	if err != nil {
		return err
	}

	_, _ = zcache.New("__account_" + h.Model.Table.Name + "__").Delete(row.Get(model.IDKey).String())

	return nil
}

func (h *AccountHandlers) CacheForID(uid interface{}) (row ztype.Map, err error) {
	idStr := ztype.ToString(uid)
	data, err := zcache.New("__account_"+h.Model.Table.Name+"__").MustGet(idStr, func(set func(data interface{},
		lifeSpan time.Duration, interval ...bool)) (err error) {
		row, err := h.Model.FindOne(func(b *builder.SelectBuilder) error {
			b.EQ(model.IDKey, uid)
			return nil
		}, false)
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
		i := user.Get(model.IDKey).Int64()
		if i == 0 {
			id = user.Get(model.IDKey).String()
		} else {
			id, err = hashid.EncryptID(h.hashid, i)
			if err != nil {
				return "", err
			}
		}
	} else {
		id = user.Get(model.IDKey).String()
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
		return nil, error_code.Unauthorized.Text("请先登录")
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

func (h *AccountHandlers) ResetManageToken(c *znet.Context, user ztype.Map, key string, expire int) {
	token, err := h.CreateManageToken(user, key, expire)
	if err == nil {
		c.SetHeader("Re-Token", token)
	}
}

func (h *AccountHandlers) QueryRoles(j *jwt.JwtInfo) (uid string, roles []string, err error) {
	if h.Model.Options.CryptID {
		var id int64
		id, err = hashid.DecryptID(h.hashid, j.U[8:])
		if err != nil {
			return "", nil, err
		}
		uid = ztype.ToString(id)
	} else {
		uid = j.U[8:]
	}

	user, err := h.CacheForID(uid)
	if err != nil {
		return "", nil, err
	}

	roles = user.Get("roles").Slice().String()
	return
}
