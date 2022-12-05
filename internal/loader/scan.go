package loader

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sohaha/zlsgo/zfile"
)

type FileType uint8

func (t FileType) Dir() string {
	switch t {
	case Model:
		return "models"
	case Flow:
		return "flows"
	case View:
		return "views"
	}
	return ""
}

func (t FileType) Suffix() string {
	switch t {
	case Model:
		return ".model.json"
	case Flow:
		return ".flow.json"
	case View:
		return ".view.json"
	}
	return ""
}

const (
	Model FileType = iota + 1
	Flow
	View
)

func Scan(root string, suffix string, recurve ...bool) (files map[string]string, dir string) {
	files = make(map[string]string)
	root = zfile.RealPath(root, true)

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && !(len(recurve) > 0 && recurve[0]) && (zfile.RealPath(path, true) != root) {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, suffix) {
			baseName, _ := filepath.Rel(root, path)
			files[strings.Replace(strings.TrimSuffix(baseName, suffix), "/", "-", -1)] = path
		}
		return err
	})

	// moduleDir := root + "module"
	// _ = filepath.WalkDir(moduleDir, func(path string, d fs.DirEntry, err error) error {
	// 	if path == moduleDir {
	// 		return nil
	// 	}

	// 	prefix := strings.Trim(strings.TrimPrefix(path, moduleDir), "/")
	// 	if prefix == "" {
	// 		return nil
	// 	}

	// 	baseName := filepath.Base(path)
	// 	if strings.HasSuffix(baseName, suffix) {
	// 		files[strings.Replace(strings.TrimSuffix(prefix, suffix), "/", "-", -1)] = path
	// 	}
	// 	return err
	// })

	return
}
