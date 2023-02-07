package job

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
	"time"
)

// MonitorClientJob 检查超时离线的客户端
func MonitorClientJob(context *tasks.TaskContext) {
	clientRepository := container.Resolve[client.Repository]()
	lst := clientRepository.ToList()
	// 检查所有客户端
	for _, do := range lst.ToArray() {
		// 检查客户端是否存活
		do.CheckOnline()
		clientRepository.Save(&do)

		// 活动时间超过5分钟的，则移除列表
		if do.IsOffline() && time.Now().Sub(do.ActivateAt).Seconds() > 300 {
			clientRepository.RemoveClient(do.Id)
		}
	}
}
