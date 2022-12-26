package app

import (
	"zlsapp/internal/account"
	"zlsapp/service"
)

func InitTasks() []service.Task {
	return []service.Task{
		account.ClearLogs(),
	}
}
