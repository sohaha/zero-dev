package loader_test

import (
	"testing"

	"zlsapp/internal/loader"
)

func TestScan(t *testing.T) {
	files := loader.Scan("../../app/models", loader.Model.Suffix())
	t.Log(files)
}
