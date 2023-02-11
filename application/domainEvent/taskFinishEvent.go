package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
)

// TaskFinishEvent 任务完成事件
func TaskFinishEvent(message any, _ core.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	// 先保存任务内容
	taskGroupRepository.SaveTask(do.Task)
	// 成功才要计算下一个周期
	if do.Task.Status == enum.Success {
		do.CalculateNextAtByCron()
	}
	// 任务初始化
	do.CreateTask()
	taskGroupRepository.SaveAndTask(*do.DomainObject)
}
