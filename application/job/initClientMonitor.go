package job

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/trace"
)

// InitClientMonitor 初始化客户端
func InitClientMonitor() {
	// 链路追踪
	traceContext := container.Resolve[trace.IManager]().EntryTask("初始化客户端")
	defer traceContext.End()

	clientRepository := container.Resolve[client.Repository]()
	// 检查所有客户端
	clientRepository.ToList().Foreach(func(clientDO *client.DomainObject) {
		domain.MonitorClientPush(clientDO)
	})
}
