package mm_test

import (
	"fmt"
	"testing"
	"zlsapp/internal/mm"
	"zlsapp/internal/parse"

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

	_, err = initModel(false)
	if err != nil {
		fmt.Println(err)
		return
	}
	m.Run()

	zfile.Rmdir("./testdata/test.db")
}

func TestModel(t *testing.T) {
	tt := zlsgo.NewTest(t)

	m2, err := initModel(false)
	tt.Equal(mm.ErrModuleAlreadyExists, err)
	zlog.Debug(err)

	m3, err := initModel(true)
	tt.NoError(err)
	zlog.Debug(err)

	_ = m2
	_ = m3

	rows, err := m3.FindAll(nil, parse.StorageOptions{
		Fields:  []string{"id"},
		OrderBy: map[string]int8{"id": -1},
	})
	tt.NoError(err)

	l := len(rows)
	tt.EqualTrue(l > 0)

	t.Log(l)
	v := rows[0].Get("id").Int()
	for i := 1; i < l-1; i++ {
		tt.EqualTrue(rows[i].Get("id").Int() < v)
		v = rows[i].Get("id").Int()
	}
}

func initDB() (*zdb.DB, error) {
	return zdb.New(&sqlite3.Config{
		File: "./testdata/test.db",
	})
}

func initModel(force bool) (*mm.Model, error) {
	m, err := mm.Add("news", modelJSON, func(m *mm.Model) (parse.Storageer, error) {
		s := parse.New(DB, m.Table.Name)
		return s, nil
	}, force)
	if err == nil {
		err = m.Migration(false).Auto()
	}
	return m, err
}