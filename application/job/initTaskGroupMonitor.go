package job

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
)

// InitTaskGroupMonitor 初始化任务组监听
func InitTaskGroupMonitor() {
	repository := container.Resolve[taskGroup.Repository]()
	repository.ToList().Foreach(func(taskGroupDO *taskGroup.DomainObject) {
		domain.MonitorTaskGroupPush(taskGroupDO)
	})
}
