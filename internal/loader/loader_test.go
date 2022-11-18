package loader_test

import (
	"testing"

	"zlsapp/internal/loader"
)

func TestScan(t *testing.T) {
	files, dir := loader.Scan("../../app", loader.Model)
	t.Log(dir, files)
}
