package jet

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/httpfs"
	"github.com/sohaha/zlsgo/zarray"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
)

type Engine struct {
	directory  string
	fileSystem http.FileSystem
	loaded     bool
	mutex      sync.RWMutex
	funcmap    map[string]interface{}
	Templates  *jet.Set
	options    Options
}

var extensions = []string{".html.jet", ".jet.html", ".jet"}

// New returns a Jet render engine for Fiber
func New(directory string, opt ...func(o *Options)) *Engine {
	o := getOption(opt...)
	if !zarray.Contains(extensions, o.Extension) {
		Log.Fatalf("%s extension is not a valid jet engine ['.html.jet', .jet.html', '.jet']", o.Extension)
	}

	engine := &Engine{
		directory: zfile.RealPath(directory),
		funcmap:   make(map[string]interface{}),
		options:   o,
	}

	return engine
}

func NewFileSystem(fs http.FileSystem, opt ...func(o *Options)) *Engine {
	o := getOption(opt...)
	if !zarray.Contains(extensions, o.Extension) {
		Log.Fatalf("%s extension is not a valid jet engine ['.html.jet', .jet.html', '.jet']", o.Extension)
	}

	engine := &Engine{
		directory:  "/",
		fileSystem: fs,
		funcmap:    make(map[string]interface{}),
		options:    o,
	}

	return engine
}

// AddFunc adds the function to the template's function map
func (e *Engine) AddFunc(name string, fn interface{}) *Engine {
	e.mutex.Lock()
	e.funcmap[name] = fn
	e.mutex.Unlock()
	return e
}

// Parse parses the templates to the engine
func (e *Engine) Load() (err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var loader jet.Loader

	if e.fileSystem != nil {
		loader, err = httpfs.NewLoader(e.fileSystem)
		if err != nil {
			return
		}
	} else {
		loader = jet.NewInMemLoader()
	}

	opts := []jet.Option{jet.WithDelims(e.options.Delims.Left, e.options.Delims.Right)}

	if e.options.Debug {
		opts = append(opts, jet.InDevelopmentMode())

	}

	e.Templates = jet.NewSet(
		loader,
		opts...,
	)

	for name, fn := range e.funcmap {
		e.Templates.AddGlobal(name, fn)
	}

	e.loaded = true

	if _, ok := loader.(*jet.InMemLoader); ok {
		total := 0
		tip := zstring.Buffer()
		err = filepath.Walk(e.directory, func(path string, info os.FileInfo, err error) error {
			l := loader.(*jet.InMemLoader)
			if err != nil {
				return err
			}
			if info == nil || info.IsDir() {
				return nil
			}
			if len(e.options.Extension) >= len(path) || path[len(path)-len(e.options.Extension):] != e.options.Extension {
				return nil
			}
			rel, err := filepath.Rel(e.directory, path)
			if err != nil {
				return err
			}
			name := strings.TrimSuffix(rel, e.options.Extension)
			buf, err := ReadFile(path, e.fileSystem)
			if err != nil {
				return err
			}

			l.Set(name, string(buf))
			if e.options.Debug {
				total++
				tip.WriteString("\t    - " + name + "\n")
			}

			return err
		})

		if err == nil && e.options.Debug {
			Log.Debugf("Loaded HTML Templates (%d): \n%s", total, tip.String())
		}
	}

	return
}

// Execute will render the template by name
func (e *Engine) Render(out io.Writer, template string, binding ztype.Map, layout ...string) error {
	if !e.loaded || e.options.Reload {
		if e.options.Reload {
			e.loaded = false
		}
		if err := e.Load(); err != nil {
			return err
		}
	}
	tmpl, err := e.Templates.GetTemplate(template)
	if err != nil || tmpl == nil {
		return fmt.Errorf("render: template %s could not be loaded: %v", template, err)
	}

	var bind jet.VarMap
	if binding != nil {
		bind = make(jet.VarMap)
		for key, value := range binding {
			bind.Set(key, value)
		}
	}
	if len(layout) > 0 && layout[0] != "" {
		lay, err := e.Templates.GetTemplate(layout[0])
		if err != nil {
			return err
		}

		bind.Set(e.options.Layout, func() {
			_ = tmpl.Execute(out, bind, empty)
		})
		return lay.Execute(out, bind, empty)
	}
	return tmpl.Execute(out, bind, nil)
}
