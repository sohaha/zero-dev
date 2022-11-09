package account

import (
	"zlsapp/internal/model"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/ztype"
	"github.com/zlsgo/zdb"
)

func migration(di zdi.Invoker) (m *model.Model, err error) {
	_, diErr := di.Invoke(func(db *zdb.DB) {
		zerror.Panic(userModel(db))
		zerror.Panic(logsModel(db))
	})

	if diErr != nil {
		return nil, diErr
	}

	m, _ = model.Get(UserModel)
	return
}

func defaultAccount() ztype.Maps {
	return ztype.Maps{
		{
			model.IDKey: 1,
			"account":   "admin",
			"password":  "admin",
			"status":    1,
			"roles":     []string{"admin"},
			"avatar":    "",
		},
	}
}
