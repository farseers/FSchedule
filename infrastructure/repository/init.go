package repository

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
)

// InitRepository 初始化仓储
func InitRepository() {
	// 注册仓储
	container.Register(func() client.Repository {
		return &clientRepository{
			Client: redis.NewClient("default"),
		}
	})

	registerTaskGroupRepository()
}
