package hashid

import (
	"zlsapp/service"

	"github.com/speps/go-hashids/v2"
)

type HashID struct {
	hash *hashids.HashID
}

func Init(conf *service.Conf) *HashID {
	salt := conf.Core().GetString("hashid.salt")
	if salt == "" {
		salt = conf.Core().GetString("account.key")
	}
	return New(salt, 8)
}

func New(salt string, MinLength int) *HashID {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = MinLength
	hash, _ := hashids.NewWithData(hd)
	return &HashID{hash: hash}
}

func EncryptID(hashI *HashID, id int64) (string, error) {
	return hashI.hash.EncodeInt64([]int64{id})
}

func DecryptID(hashI *HashID, hid string) (int64, error) {
	ids, err := hashI.hash.DecodeInt64WithError(hid)
	if err != nil {
		return 0, err
	}
	return ids[0], nil
}
