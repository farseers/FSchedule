package repository

import (
	"FSchedule/domain/client"
	"FSchedule/domain/schedule"
	"github.com/farseer-go/fs/container"
)

// InitRepository 初始化仓储
func InitRepository() {
	// 注册client仓储
	container.Register(func() client.Repository {
		return &clientRepository{}
	})
	// 注册client仓储
	container.Register(func() schedule.Repository {
		return &scheduleRepository{}
	})

	registerTaskGroupRepository()
}
