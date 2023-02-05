package job

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
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
	}
}
