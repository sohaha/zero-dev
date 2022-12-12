package account

import (
	"zlsapp/common"
	"zlsapp/internal/parse"
	"zlsapp/service"

	"github.com/mileusna/useragent"
	"github.com/sohaha/zlsgo/zjson"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
	ipRegion "github.com/zlsgo/ip"
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

type LogType int

const (
	LogTypeCommon LogType = iota + 1
	LogTypeAction
	LogTypeLogin
)

func (l *Logs) Init(z *znet.Engine) {
	l.model, _ = parse.GetModel(LogsModel)
}

func WrapActionLogs(c *znet.Context, action, module string) {
	WrapLogs(c, action, func(data *ztype.Map) {
		(*data)["category"] = LogTypeAction
		(*data)["module"] = module
		p := c.PrevContent()
		success := p.Code.Load() == 200
		(*data)["result"] = success
		if !success {
			json := zjson.ParseBytes(p.Content)
			(*data)["detail"] = json.Get("msg").String()
		}
	})
}

// var logPool = zpool.New(100)

func init() {
	go func() {
		_, _ = ipRegion.Region("")
	}()
}

func WrapLogs(c *znet.Context, action string, fn ...func(data *ztype.Map)) {
	// _ = logPool.Do(func() {
	status := LogsStatusRead
	uid := common.GetUID(c)
	url := c.Request.URL.Path
	m, _ := parse.GetModel(LogsModel)

	ip := c.GetClientIP()
	region, _ := ipRegion.Region(ip)
	u := useragent.Parse(c.GetUserAgent())

	data := ztype.Map{
		"uid":             uid,
		"action":          action,
		"category":        LogTypeCommon,
		"ip":              ip,
		"detail":          "",
		"method":          c.Request.Method,
		"path":            url,
		"os":              u.OS,
		"ip_region":       zstring.TrimSpace(region.Country + " " + region.Province + " " + region.City),
		"os_version":      u.OSVersion,
		"device":          u.Device,
		"browser":         u.Name,
		"browser_version": u.Version,
		"status":          status,
	}

	for _, f := range fn {
		f(&data)
	}

	_, _ = parse.Insert(m, data)
	// })
}
