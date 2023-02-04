package job

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
)

// MonitorJob 任务组监听
func MonitorJob() {
	repository := container.Resolve[taskGroup.Repository]()
	lst := repository.ToList()
	for _, taskGroupDO := range lst.ToArray() {
		domain.MonitorPush(taskGroupDO)
	}
}
