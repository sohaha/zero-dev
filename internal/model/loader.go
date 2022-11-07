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

func Loader(db *zdb.DB, path, name string) error {
	json, _ := zfile.ReadFile(path)

	zerror.Panic(ValidateModelSchema(json))
	// baseName := filepath.Base(path)
	m, err := Add(db, name, json)
	zlog.Debug(name)
	_ = m
	return err
}

func NewLoader(di zdi.Invoker) error {
	_, err := di.Invoke(func(db *zdb.DB) {
		filepath.WalkDir(zfile.RealPath("./app/model"), func(path string, d fs.DirEntry, err error) error {
			// zlog.Debug(path, d)
			baseName := filepath.Base(path)
			// ext := filepath.Ext(baseName)
			if strings.HasSuffix(baseName, ".model.json") {
				zlog.Debug(baseName)
				name := strings.TrimSuffix(baseName, ".model.json")
				err = Loader(db, path, name)
				if err != nil {
					zlog.Error(err)
				}
			}
			// zlog.Debug(baseName, ext)
			return err
		})
	})
	globalModels.ForEach(func(key string, value *Model) bool {
		zlog.Debug(key, value)
		value.Migration().Auto()
		return true
	})
	os.Exit(0)
	return err
}
