package loader

import (
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zstring"
)

func modelLog(v ...interface{}) {
	d := []interface{}{
		zlog.ColorTextWrap(zlog.ColorLightMagenta, zstring.Pad("Model", 6, " ", zstring.PadLeft)),
	}
	d = append(d, v...)
	zlog.Debug(d...)
}
