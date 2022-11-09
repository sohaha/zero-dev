package model

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sohaha/zlsgo/zdi"
	"github.com/sohaha/zlsgo/zerror"
	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/zlsgo/zdb"
)

func Loader(db *zdb.DB, path, name string, force bool) error {
	json, _ := zfile.ReadFile(path)

	zerror.Panic(ValidateModelSchema(json))
	// baseName := filepath.Base(path)
	m, err := Add(db, name, json, force)
	zlog.Debug(name)
	_ = m
	return err
}

func NewLoader(di zdi.Invoker) error {
	suffix := ".model.json"
	_, err := di.Invoke(func(db *zdb.DB) {
		filepath.WalkDir(zfile.RealPath("./app/model"), func(path string, d fs.DirEntry, err error) error {
			baseName := filepath.Base(path)
			if strings.HasSuffix(baseName, suffix) {
				zlog.Debug(baseName)
				name := strings.TrimSuffix(baseName, suffix)
				err = Loader(db, path, name, false)
				if err != nil {
					zlog.Error(err)
				}
			}
			// zlog.Debug(baseName, ext)
			return err
		})
	})

	os.Exit(0)
	return err
}
