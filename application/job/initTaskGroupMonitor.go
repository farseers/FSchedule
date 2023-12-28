package job

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/trace"
)

// InitTaskGroupMonitor 初始化任务组监听
func InitTaskGroupMonitor() {
	// 链路追踪
	traceContext := container.Resolve[trace.IManager]().EntryTask("初始化任务组监听")
	defer traceContext.End()

	repository := container.Resolve[taskGroup.Repository]()
	repository.ToList().Foreach(func(taskGroupDO *taskGroup.DomainObject) {
		domain.MonitorTaskGroupPush(taskGroupDO)
	})
}
