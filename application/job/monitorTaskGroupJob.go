package job

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

// MonitorTaskGroupJob 任务组监听
func MonitorTaskGroupJob(context *tasks.TaskContext) {
	repository := container.Resolve[taskGroup.Repository]()
	lst := repository.ToList()
	for _, taskGroupDO := range lst.ToArray() {
		domain.MonitorTaskGroupPush(&taskGroupDO)
	}
}
