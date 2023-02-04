package repository

import (
	"FSchedule/domain/client"
	"FSchedule/domain/schedule"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
)

// InitRepository 初始化仓储
func InitRepository() {
	// 注册client仓储
	container.Register(func() client.Repository {
		return &clientRepository{
			Client: redis.NewClient("default"),
		}
	})
	// 注册client仓储
	container.Register(func() schedule.Repository {
		return &scheduleRepository{
			Client: redis.NewClient("default"),
		}
	})

	registerTaskGroupRepository()
}
