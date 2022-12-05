package loader_test

import (
	"testing"

	"zlsapp/internal/loader"
)

func TestScan(t *testing.T) {
	files, dir := loader.Scan("../../app/models", loader.Model.Suffix())
	t.Log(dir, files)
}
