package domainEvent

import (
	"FSchedule/domain/enum/scheduleStatus"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
)

// TaskFinishEvent 任务完成事件
func TaskFinishEvent(message any, _ core.EventArgs) {
	do := message.(*taskGroup.DomainObject)
	if do.Task.ScheduleStatus != scheduleStatus.Fail && !do.Task.IsFinish() {
		return
	}

	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	// 先保存任务内容
	taskGroupRepository.SaveTask(do.Task)

	// 计算下一个周期
	if do.CalculateNextAtByCron() {
		//status := do.Task.ExecuteStatus
		// 任务初始化
		do.CreateTask()
		//flog.Debugf("任务组：%s %d 任务完成，结果：%s，下次执行时间：%s", do.Name, do.Task.Id, status.String(), do.Task.StartAt.ToString("yyyy-MM-dd HH:mm:ss"))
	}
	taskGroupRepository.SaveAndTask(*do)
}
