package migration

import (
	"zlsapp/app/model"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/zlsgo/zdb"
)

func RunMigrations(di zdi.Invoker) error {
	_, err := di.Invoke(func(db *zdb.DB) {
		name := "testdata/user.model.json"
		json, _ := zfile.ReadFile(name)

		zerror.Panic(model.ValidateModelSchema(json))

		m, err := model.Add(db, name, json)

		zerror.Panic(err)

		migration := m.Migration()

		zerror.Panic(migration.Auto())
	})

	return err
}
