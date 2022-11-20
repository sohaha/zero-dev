package account

import (
	"errors"
	"zlsapp/common"
	"zlsapp/internal/parse"

	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/ztype"
)

func WrapLogs(c *znet.Context, action, remark string, hint ...bool) error {
	status := 2
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

	data := ztype.Map{
		"uid":        uid,
		"action":     action,
		"ip":         ip,
		"remark":     remark,
		"user_agent": userAgent,
		"status":     status,
	}
	_, err := parse.Insert(m, data)
	return err
}
