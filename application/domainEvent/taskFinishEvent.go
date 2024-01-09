package domainEvent

import (
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

// TaskFinishEvent 任务完成事件
func TaskFinishEvent(message any, _ core.EventArgs) {
	do := message.(*taskGroup.DomainObject)
	if !do.Task.IsFinish() {
		return
	}

	// 链路追踪
	traceContext := container.Resolve[trace.IManager]().EntryTaskGroup("任务完成事件", do.Name, do.Task.Id)
	defer traceContext.End()

	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	// 先保存任务内容
	taskGroupRepository.SaveTask(do.Task)
	// 计算下一个周期
	if do.CalculateNextAtByCron() {
		// 任务初始化
		do.CreateTask()
		flog.Debugf("任务组：%s %d 任务完成，下次执行时间：%s", do.Name, do.Task.Id, do.Task.StartAt.ToString("yyyy-MM-dd HH:mm:ss"))
	}
	taskGroupRepository.SaveAndTask(*do)
}
