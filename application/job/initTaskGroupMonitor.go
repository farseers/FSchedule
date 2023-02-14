package job

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
)

// InitTaskGroupMonitor 初始化任务组监听
func InitTaskGroupMonitor() {
	repository := container.Resolve[taskGroup.Repository]()
	lst := repository.ToList()
	for i := 0; i < lst.Count(); i++ {
		taskGroupDO := lst.Index(i)
		domain.MonitorTaskGroupPush(&taskGroupDO)
	}
}
