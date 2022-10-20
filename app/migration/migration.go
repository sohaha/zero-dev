package migration

import (
	"zlsapp/app/model"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/zdb"
)

func RunMigrations(di zdi.Invoker) error {
	_, err := di.Invoke(func(db *zdb.DB) {
		json, _ := zfile.ReadFile("testdata/user.model.json")
		m, err := model.ParseJSON(db, json)

		zerror.Panic(err)

		migration := m.Migration()

		exist := migration.HasTable()
		zlog.Debug(exist)
		if !exist {
			zerror.Panic(migration.CreateTable())
			zlog.Debug("新建")
		} else {
			zerror.Panic(migration.UpdateTable())
			zlog.Debug("需要更新表结构")
		}

		zlog.Debug(migration.InitValue())
		// table := builder.NewTable("user").Create()

		// // schema.NewField()
		// table.Column()
		// zlog.Debug(table.Build())
	})

	return err
}
