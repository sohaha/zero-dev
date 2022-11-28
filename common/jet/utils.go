package jet

import (
	"io"
	"net/http"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zlog"
)

var empty = struct{}{}

var Log = zlog.New("jet ")

func init() {
	Log.ResetFlags(zlog.BitLevel)
}

type Delims struct {
	Left  string
	Right string
}

type Options struct {
	Extension string
	Layout    string
	Debug     bool
	Reload    bool
	Delims    Delims
}

func getOption(opt ...func(o *Options)) Options {
	o := Options{
		Extension: ".jet.html",
		Delims: Delims{
			Left:  "{{",
			Right: "}}",
		},
		Layout: "slot",
	}
	for _, f := range opt {
		f(&o)
	}
	return o
}

func ReadFile(path string, fs http.FileSystem) ([]byte, error) {
	if fs != nil {
		file, err := fs.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return io.ReadAll(file)
	}
	return zfile.ReadFile(path)
}
