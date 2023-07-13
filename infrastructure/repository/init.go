package repository

import (
	"FSchedule/domain/client"
	"FSchedule/domain/schedule"
	"FSchedule/domain/serverNode"
	"FSchedule/domain/taskLog"
	"github.com/farseer-go/fs/container"
)

// InitRepository 初始化仓储
func InitRepository() {
	// 注册serverNode仓储
	container.Register(func() serverNode.Repository {
		return &serverNodeRepository{}
	})

	// 注册client仓储
	container.Register(func() client.Repository {
		return &clientRepository{}
	})

	// 注册schedule仓储
	container.Register(func() schedule.Repository {
		return &scheduleRepository{}
	})

	// 注册taskLog仓储
	container.Register(func() taskLog.Repository {
		return &TaskLogRepository{}
	})

	// 注册taskGroup仓储
	registerTaskGroupRepository()
}
