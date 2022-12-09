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

func Scan(root string, suffix string, recurve ...bool) (files []string) {
	files = make([]string, 0)
	root = zfile.RealPath(root, true)

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		path = zfile.RealPath(path)
		if err != nil {
			return err
		}

		if d.IsDir() && !(len(recurve) > 0 && recurve[0]) && (zfile.RealPath(path, true) != root) {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, suffix) {
			// baseName := zfile.SafePath(path, root)
			// baseName = strings.Replace(baseName, "\\", "/", -1)
			// files[strings.Replace(strings.TrimSuffix(baseName, suffix), "/", "-", -1)] = path
			files = append(files, path)
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

func toName(path string, root string) string {
	baseName := zfile.SafePath(zfile.RealPath(path), zfile.RealPath(root))
	sp := strings.Split(baseName, ".")
	if len(sp) > 2 {
		baseName = strings.Join(sp[:len(sp)-2], ".")
	}
	return strings.Replace(baseName, "/", "-", -1)
}
