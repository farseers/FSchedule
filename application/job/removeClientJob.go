package job

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/tasks"
)

// RemoveClientJob 20分钟检查客户端是否永久离线
func RemoveClientJob(context *tasks.TaskContext) {
	clientRepository := container.Resolve[client.Repository]()
	// 检查所有客户端
	clientRepository.ToList().Foreach(func(clientDO *client.DomainObject) {
		if clientDO.IsOffline() && dateTime.Since(clientDO.ActivateAt).Minutes() >= 10 {
			clientRepository.RemoveClient(clientDO.Id)
		}
	})
}
