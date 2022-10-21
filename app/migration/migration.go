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
		name := "testdata/user.model.json"
		json, _ := zfile.ReadFile(name)
		m, err := model.Add(db, name, json)

		zerror.Panic(err)

		migration := m.Migration()

		exist := migration.HasTable()
		if !exist {
			zlog.Debug("新建")
			zerror.Panic(migration.CreateTable())
		} else {
			zlog.Debug("需要更新表结构")
			zerror.Panic(migration.UpdateTable())
		}

		zlog.Debug("初始化数据")

		err = migration.InitValue()
		zerror.Panic(err)
		// table := builder.NewTable("user").Create()

		// // schema.NewField()
		// table.Column()
		// zlog.Debug(table.Build())
	})

	return err
}
