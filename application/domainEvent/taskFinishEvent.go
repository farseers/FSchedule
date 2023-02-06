package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/container"
)

func TaskFinishEvent(message any, _ eventBus.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	taskGroupDO := taskGroupRepository.ToEntity(do.Name)
	// 任务初始化
	taskGroupDO.CreateTask()
	taskGroupRepository.SaveAndTask(taskGroupDO)
}
