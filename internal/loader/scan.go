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
	}
	return ""
}

func (t FileType) Suffix() string {
	switch t {
	case Model:
		return ".model.json"
	case Flow:
		return ".flow.json"
	}
	return ""
}

const (
	Model FileType = iota + 1
	Flow
)

func Scan(root string, filetype FileType) (files map[string]string, dir string) {
	files = make(map[string]string)
	root = zfile.RealPath(root, true)
	suffix := filetype.Suffix()

	dir = zfile.RealPath(root + filetype.Dir())
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && (zfile.RealPath(path) != dir) {
			return filepath.SkipDir
		}
		baseName := filepath.Base(path)
		if strings.HasSuffix(baseName, suffix) {
			files[strings.TrimSuffix(baseName, suffix)] = path
		}
		return err
	})

	moduleDir := root + "module"
	_ = filepath.WalkDir(moduleDir, func(path string, d fs.DirEntry, err error) error {
		if path == moduleDir {
			return nil
		}

		prefix := strings.Trim(strings.TrimPrefix(path, moduleDir), "/")
		if prefix == "" {
			return nil
		}

		baseName := filepath.Base(path)
		if strings.HasSuffix(baseName, suffix) {
			files[strings.TrimSuffix(prefix, suffix)] = path
		}
		return err
	})

	return
}
