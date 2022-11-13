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
			"account":   "manage",
			"password":  "123456",
			"status":    1,
			"nickname":  "管理员",
			"roles":     []string{"admin"},
			"avatar":    "data:image/svg+xml,%3Csvg viewBox='0 0 36 36' fill='none' role='img' xmlns='http://www.w3.org/2000/svg' width='128' height='128'%3E%3Ctitle%3EMary Roebling%3C/title%3E%3Cmask id='mask__beam' maskUnits='userSpaceOnUse' x='0' y='0' width='36' height='36'%3E%3Crect width='36' height='36' fill='%23FFFFFF'%3E%3C/rect%3E%3C/mask%3E%3Cg mask='url(%23mask__beam)'%3E%3Crect width='36' height='36' fill='%23f0f0d8'%3E%3C/rect%3E%3Crect x='0' y='0' width='36' height='36' transform='translate(5 -1) rotate(155 18 18) scale(1.2)' fill='%23000000' rx='6'%3E%3C/rect%3E%3Cg transform='translate(3 -4) rotate(-5 18 18)'%3E%3Cpath d='M15 21c2 1 4 1 6 0' stroke='%23FFFFFF' fill='none' stroke-linecap='round'%3E%3C/path%3E%3Crect x='14' y='14' width='1.5' height='2' rx='1' stroke='none' fill='%23FFFFFF'%3E%3C/rect%3E%3Crect x='20' y='14' width='1.5' height='2' rx='1' stroke='none' fill='%23FFFFFF'%3E%3C/rect%3E%3C/g%3E%3C/g%3E%3C/svg%3E",
		},
	}
}
