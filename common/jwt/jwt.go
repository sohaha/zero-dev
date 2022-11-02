package jwt

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
)

type JwtInfo struct {
	U string `json:"u"`
	jwt.StandardClaims
}

func ParsingToken(token string, tokenKey string) (*JwtInfo, error) {
	t, err := jwt.ParseWithClaims(token, &JwtInfo{}, func(token *jwt.Token) (i interface{}, err error) {
		return zstring.String2Bytes(tokenKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*JwtInfo); ok && t.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetToken(c *znet.Context) string {
	authorization := c.GetHeader("Authorization")
	slen := len("Basic ")
	if len(authorization) > slen {
		authorization = zstring.TrimSpace(authorization[slen:])
		split := strings.Split(authorization, ".")
		if len(split) == 3 {
			return authorization
		}
		v, err := zstring.Base64Decode(zstring.String2Bytes(authorization))
		if err != nil {
			return ""
		}
		return strings.Split(zstring.Bytes2String(v), ":")[0]
	}
	return c.DefaultFormOrQuery("Authorization", "")
}
