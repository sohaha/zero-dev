package parse_test

import (
	"fmt"
	"testing"
	"zlsapp/internal/parse"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
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

	_, err := initModel(false)
	tt.Equal(parse.ErrModuleAlreadyExists, err)
	t.Log(err)

	m, err := initModel(true)
	tt.NoError(err)
	t.Log(err)

	id, err := parse.Insert(m, map[string]interface{}{
		"title":    "test",
		"key":      "123",
		"category": 1,
		"content":  "test content",
	})
	t.Log(id)
	tt.NoError(err)

	row, err := parse.FindOne(m, ztype.Map{
		parse.IDKey: id,
	}, func(so *parse.StorageOptions) error {
		so.Fields = []string{"key", "title"}
		return nil
	})
	t.Log(row)
	tt.NoError(err)
	tt.Equal(zstring.Md5("123"), row.Get("key").String())
	tt.Equal(32, len(row.Get("key").String()))
	tt.Equal("", row.Get("content").String())
	tt.Equal("test", row.Get("title").String())

	row, err = parse.FindOne(m, ztype.Map{
		parse.IDKey + " >": 1,
		parse.IDKey:        3,
		"":                 "reading != 11",
	})
	t.Log(row)
	tt.NoError(err)
	tt.Equal(3, row.Get(parse.IDKey).Int())

	rows, err := parse.Find(m, ztype.Map{}, func(so *parse.StorageOptions) error {
		so.Fields = []string{parse.IDKey}
		so.OrderBy = map[string]int8{parse.IDKey: -1}
		return nil
	})
	tt.NoError(err)

	l := len(rows)
	tt.EqualTrue(l > 0)
	t.Log(l)
	v := rows[0].Get(parse.IDKey).Int()
	for i := 1; i < l-1; i++ {
		tt.EqualTrue(rows[i].Get("id").Int() < v)
		v = rows[i].Get(parse.IDKey).Int()
	}

	// _, err = m.Delete()
}

func initDB() (*zdb.DB, error) {
	return zdb.New(&sqlite3.Config{
		File: "./testdata/test.db",
	})
}

func initModel(force bool) (*parse.Model, error) {
	m, err := parse.AddModel("news", modelJSON, func(m *parse.Model) (parse.Storageer, error) {
		s := parse.NewSQL(DB, m.Table.Name)
		return s, nil
	}, force)
	if err == nil {
		err = m.Migration().Auto(false)
	}
	return m, err
}
