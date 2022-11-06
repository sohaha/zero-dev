package conf

import "github.com/sohaha/zlsgo/zstring"

type Account struct {
	Key    string
	IDKey  string `mapstructure:"idKey"`
	Expire int
	Only   bool
}

func init() {
	DefaultConf = append(DefaultConf, Account{
		Key:    zstring.Rand(8),
		IDKey:  "",
		Expire: 0,
		Only:   false,
	})
}
