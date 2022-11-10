package open

import (
	"zlsapp/service"

	"github.com/sohaha/zlsgo/znet"
)

type Open struct {
	service.App
	Path string
}

func (h *Open) Init(r *znet.Engine) {

}
