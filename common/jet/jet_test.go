package jet

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/sohaha/zlsgo"
	"github.com/sohaha/zlsgo/zfile"
)

func trim(str string) string {
	trimmed := strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(str, " "))
	trimmed = strings.Replace(trimmed, " <", "<", -1)
	trimmed = strings.Replace(trimmed, "> ", ">", -1)
	return trimmed
}

func TestRender(t *testing.T) {
	tt := zlsgo.NewTest(t)

	engine := New("./views", func(o *Options) {
		o.Debug = true
		o.Extension = ".jet.html"
	})

	err := engine.Load()
	tt.NoError(err)

	var buf bytes.Buffer
	err = engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Hello, World!",
	})
	tt.NoError(err)
	expect := `<h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2>`

	tt.Equal(expect, trim(buf.String()))

	buf.Reset()

	err = engine.Render(&buf, "errors/404", map[string]interface{}{
		"Title": "Hello, World!",
	})
	tt.NoError(err)
	expect = `<h1>Hello, World!</h1>`
	tt.Equal(expect, trim(buf.String()))
}

func TestLayout(t *testing.T) {
	tt := zlsgo.NewTest(t)

	engine := New("./views")

	err := engine.Load()
	tt.NoError(err)

	var buf bytes.Buffer
	err = engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Hello, World!",
	}, "layouts/main")
	tt.NoError(err)

	expect := `<!DOCTYPE html><html><head><title>Title</title></head><body><h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2></body></html>`
	tt.Equal(expect, trim(buf.String()))
}

func TestEmptyLayout(t *testing.T) {
	tt := zlsgo.NewTest(t)
	engine := New("./views")

	var buf bytes.Buffer

	err := engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Hello, World!",
	}, "")
	tt.NoError(err)
	expect := `<h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2>`
	tt.Equal(expect, trim(buf.String()))
}

func TestFileSystem(t *testing.T) {
	tt := zlsgo.NewTest(t)
	engine := NewFileSystem(http.Dir(zfile.RealPath("./views")), func(o *Options) {
		o.Debug = true
	})

	var buf bytes.Buffer
	err := engine.Render(&buf, "index", map[string]interface{}{
		"Title": "Hello, World!",
	}, "/layouts/main")
	tt.NoError(err)

	expect := `<!DOCTYPE html><html><head><title>Title</title></head><body><h2>Header</h2><h1>Hello, World!</h1><h2>Footer</h2></body></html>`
	tt.Equal(expect, trim(buf.String()))
}

func TestReload(t *testing.T) {
	tt := zlsgo.NewTest(t)
	engine := NewFileSystem(http.Dir("./views"), func(o *Options) {
		o.Reload = true
	})

	err := engine.Load()
	tt.NoError(err)

	err = zfile.WriteFile("./views/reload.jet.html", []byte("after reload\n"))
	tt.NoError(err)

	defer func() {
		_ = zfile.WriteFile("./views/reload.jet.html", []byte("before reload\n"))
	}()

	_ = engine.Load()

	var buf bytes.Buffer
	err = engine.Render(&buf, "reload", nil)

	tt.NoError(err)
	expect := "after reload"
	tt.Equal(expect, trim(buf.String()))
}
