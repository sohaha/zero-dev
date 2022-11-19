package model_test

import (
	"fmt"
	"testing"
	"zlsapp/internal/model"
	"zlsapp/internal/model/storage"
	"zlsapp/internal/model/storage/sql"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/zdb"
	"github.com/zlsgo/zdb/driver/sqlite3"
)

var DB *zdb.DB
var modelJSON, _ = zfile.ReadFile("./testdata/news.model.json")

func TestMain(m *testing.M) {
	db, err := initDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	DB = db

	_, err = getModel(false)
	if err != nil {
		fmt.Println(err)
		return
	}
	m.Run()
}

func initDB() (*zdb.DB, error) {
	return zdb.New(&sqlite3.Config{
		File: "./testdata/test.db",
	})
}

func TestModel(t *testing.T) {
	tt := zlsgo.NewTest(t)

	m2, err := getModel(false)
	tt.Equal(model.ErrModuleAlreadyExists, err)
	zlog.Debug(err)

	m3, err := getModel(true)
	tt.NoError(err)
	zlog.Debug(err)

	_ = m2
	_ = m3
}

func getModel(force bool) (*model.Model, error) {
	m, err := model.Add("news", modelJSON, func(m *model.Model) (storage.Storageer, error) {
		s := sql.New(DB, m.Table.Name)
		return s, nil
	}, force)
	if err == nil {
		// err = m.Migration().Auto(false)
	}
	return m, err
}
