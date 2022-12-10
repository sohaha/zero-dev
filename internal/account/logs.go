package account

import (
	"errors"
	"zlsapp/common"
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/mileusna/useragent"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

type Logs struct {
	service.App
	Path  string
	model *parse.Modeler
}

const (
	LogsStatusUnread = iota + 1
	LogsStatusRead
)

func (l *Logs) Init(z *znet.Engine) {
	l.model, _ = parse.GetModel(LogsModel)
}

func WrapLogs(c *znet.Context, action, remark string, hint ...bool) error {
	status := LogsStatusRead
	if len(hint) > 0 && hint[0] {
		status = 1
	}
	uid := common.GetUID(c)
	return CreateLogs(uid, action, remark, c.GetClientIP(), c.GetUserAgent(), status)
}

func CreateLogs(uid, action, remark, ip, userAgent string, status int) error {
	m, ok := parse.GetModel(LogsModel)
	if !ok {
		return errors.New("model(" + LogsModel + ") not found")
	}

	u := useragent.Parse(userAgent)
	data := ztype.Map{
		"uid":        uid,
		"action":     action,
		"ip":         ip,
		"remark":     remark,
		"os":         u.OS,
		"os_version": u.OSVersion,
		"device":     u.Device,
		// "user_agent": userAgent,
		"status": status,
	}
	_, err := parse.Insert(m, data)
	return err
}
