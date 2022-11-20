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

	_, err = initModel(true)
	tt.NoError(err)
	t.Log(err)
}

func TestModelAction(t *testing.T) {
	tt := zlsgo.NewTest(t)

	m, err := initModel(true)
	tt.NoError(err)

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
	tt.Equal(2, len(row))
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

	row, err = parse.FindOne(m, 2)
	t.Log(row)
	tt.NoError(err)

	update := ztype.Map{
		"title":    "第 2 篇",
		"category": 2,
		"content":  "替换新闻",
		"key":      "new",
	}

	total, err := parse.Update(m, "2", update)
	tt.NoError(err)
	tt.Equal(int64(1), total)

	row2, err := parse.FindOne(m, ztype.Map{
		"": "title = '第 2 篇'",
	})
	t.Log(row2)
	tt.NoError(err)

	tt.EqualTrue(row2.Get("title").String() != row.Get("title").String())
	tt.EqualTrue(row2.Get("category").String() == row.Get("category").String())
	tt.EqualTrue(row2.Get("content").String() != row.Get("content").String())
	tt.EqualTrue(row2.Get("key").String() == row.Get("key").String())

	items, pages, err := parse.Pages(m, 2, 2, ztype.Map{}, func(so *parse.StorageOptions) error {
		so.Fields = []string{parse.IDKey}
		return nil
	})
	tt.NoError(err)
	t.Log(items)
	t.Log(pages)

	total, err = parse.Delete(m, 2)
	tt.NoError(err)
	tt.Equal(int64(1), total)

	row, err = parse.FindOne(m, 2)
	t.Log(row)

	tt.Equal(zdb.ErrNotFound, err)

}

func initDB() (*zdb.DB, error) {
	return zdb.New(&sqlite3.Config{
		File: "./testdata/test.db",
	})
}

func initModel(force bool) (*parse.Modeler, error) {
	m, err := parse.AddModel("news", modelJSON, func(m *parse.Modeler) (parse.Storageer, error) {
		s := parse.NewSQL(DB, m.Table.Name)
		return s, nil
	}, force)
	if err == nil {
		err = m.Migration().Auto(false)
	}
	return m, err
}
