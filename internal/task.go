package app

import (
	"zlsapp/internal/account"
	"zlsapp/service"
)

func InitTask() []service.Task {
	return []service.Task{
		account.ClearLogs(),
	}
}
