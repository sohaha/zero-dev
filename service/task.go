package service

import (
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztime/cron"
)

type Task struct {
	Name string
	Cron string
	Run  func(app *App)
}

func RunTask(tasks []Task, app *App) (err error) {
	t := cron.New()

	log := func(v ...interface{}) {
		d := []interface{}{
			zlog.ColorTextWrap(zlog.ColorLightCyan, zstring.Pad("Cron", 6, " ", zstring.PadLeft)),
		}
		d = append(d, v...)
		zlog.Debug(d...)
	}

	for _, task := range tasks {
		_, err = t.Add(task.Cron, func() {
			task.Run(app)
		})

		if err != nil {
			return
		}

		next, _ := cron.ParseNextTime(task.Cron)
		log("Register: " + zlog.Log.ColorTextWrap(zlog.ColorLightGreen, task.Name) + zlog.ColorTextWrap(zlog.ColorLightWhite, " ["+task.Cron+"] -> ["+ztime.FormatTime(next)+"]"))
	}

	t.Run()
	return nil
}
