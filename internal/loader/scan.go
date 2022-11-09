package loader

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sohaha/zlsgo/zfile"
)

type FileType uint8

func (t FileType) String() string {
	switch t {
	case Model:
		return "model"
	}
	return ""
}

const (
	Model FileType = iota + 1
)

func Scan(root string, filetype FileType) map[string]string {
	files := make(map[string]string)
	root = zfile.RealPath(root, true)
	suffix := "." + filetype.String() + ".json"

	dir := zfile.RealPath(root + filetype.String())
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

	return files
}
