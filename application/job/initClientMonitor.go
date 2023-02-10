package job

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/container"
)

// InitClientMonitor 初始化客户端
func InitClientMonitor() {
	clientRepository := container.Resolve[client.Repository]()
	lst := clientRepository.ToList()
	// 检查所有客户端
	for _, clientDO := range lst.ToArray() {
		domain.MonitorClientPush(&clientDO)
	}
}
