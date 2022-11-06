package hashid

import (
	"zlsapp/service"

	"github.com/speps/go-hashids/v2"
)

func Init(conf *service.Conf) *hashids.HashID {
	hd := hashids.NewData()
	salt := conf.Core().GetString("hashid.salt")
	if salt == "" {
		salt = conf.Core().GetString("account.key")
	}
	hd.Salt = salt
	hd.MinLength = 5
	hashI, _ := hashids.NewWithData(hd)
	return hashI
}

func EncryptID(hashI *hashids.HashID, id int64) (string, error) {
	return hashI.EncodeInt64([]int64{id})
}

func DecryptID(hashI *hashids.HashID, hid string) (int64, error) {
	ids, err := hashI.DecodeInt64WithError(hid)
	if err != nil {
		return 0, err
	}
	return ids[0], nil
}
