package parse

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zstring"
	"golang.org/x/crypto/bcrypt"
)

type cryptProcess func(string) (string, error)

func (m *Modeler) EncodeID(id any) {

}

func (m *Modeler) DecodeID(id any) {
	// v := ztype.New(id)
	// i := v.Int64()
	//
	//			if i == 0 {
	//				id = v.String()
	//			} else {
	//					_, _ = h.Di.Invoke(func(hashid *hashids.HashID) {
	//			h.Handlers = &AccountHandlers{
	//				Model:  h.Model,
	//				hashid: hashid,
	//			}
	//		})
	//				id, err = hashid.EncryptID(h.hashid, i)
	//				if err != nil {
	//					return "", err
	//				}
	//			}
	//	 hashid.EncryptID(h.hashid, i)
}

func (m *Modeler) GetCryptProcess(cryptName string) (fn cryptProcess, err error) {
	switch strings.ToLower(cryptName) {
	default:
		return nil, errors.New("crypt name not found")
	case "md5":
		fn = func(s string) (string, error) {
			return zstring.Md5(s), nil
		}
	case "password":
		fn = func(s string) (string, error) {
			bcost := bcrypt.DefaultCost
			bytes, err := bcrypt.GenerateFromPassword(zstring.String2Bytes(s), bcost)
			if err != nil {
				return "", err
			}
			return zstring.Bytes2String(bytes), nil
		}
	}
	return
}
