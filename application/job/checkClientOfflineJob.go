package job

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

// CheckClientOfflineJob 检查超时离线的客户端
func CheckClientOfflineJob(context *tasks.TaskContext) {
	clientRepository := container.Resolve[client.Repository]()
	lst := clientRepository.ToList()
	// 检查所有客户端
	for _, do := range lst.ToArray() {
		// 检查客户端是否存活
		do.CheckOnline()

		// 非离线状态，保存信息
		if do.Status != enum.Offline {
			clientRepository.Save(do)
		}
	}
}
