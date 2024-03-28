package domainEvent

import (
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
)

// TaskFinishEvent 任务完成事件
func TaskFinishEvent(message any, _ core.EventArgs) {
	do := message.(*taskGroup.DomainObject)
	if !do.Task.IsFinish() {
		return
	}

	status := do.Task.ExecuteStatus
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	// 先保存任务内容
	taskGroupRepository.SaveTask(do.Task)

	// 计算下一个周期
	if do.CalculateNextAtByCron() {
		// 任务初始化
		do.CreateTask()
		flog.Debugf("任务组：%s %d 任务完成，结果：%s，下次执行时间：%s", do.Name, do.Task.Id, status.String(), do.Task.StartAt.ToString("yyyy-MM-dd HH:mm:ss"))
	}
	taskGroupRepository.SaveAndTask(*do)
}
