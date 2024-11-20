package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"

	"github.com/bytedance/sonic"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/trace"
)

// TaskGroupUpdateSubscribe 任务组有更新（Redis订阅）
func TaskGroupUpdateSubscribe(message any, _ core.EventArgs) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	var taskGroupDO taskGroup.DomainObject
	err := sonic.Unmarshal([]byte(message.(string)), &taskGroupDO)
	if err != nil {
		return
	}

	// 通知处理该任务组的服务端，需要调用客户端发起Kill请求
	domain.GetTaskGroupMonitorByName(taskGroupDO.Name).Foreach(func(item **domain.TaskGroupMonitor) {
		taskGroupMonitor := *item
		// 之前是运行状态，改为停止状态，则需要退出调度线程
		if taskGroupMonitor.IsEnable && !taskGroupDO.IsEnable {
			// 主动通知客户端，停止任务
			taskGroupMonitor.TaskKill()
		}
		taskGroupMonitor.DomainObject.Data = taskGroupDO.Data
		taskGroupMonitor.DomainObject.Caption = taskGroupDO.Caption
		taskGroupMonitor.DomainObject.StartAt = taskGroupDO.StartAt

		// 手动修改了执行时间
		if taskGroupMonitor.DomainObject.NextAt.ToString("yyyy-MM-dd HH:mm:ss") != taskGroupDO.NextAt.ToString("yyyy-MM-dd HH:mm:ss") {
			taskGroupMonitor.DomainObject.NextAt = taskGroupDO.NextAt
			taskGroupMonitor.DomainObject.Task.StartAt = taskGroupDO.NextAt
		}
		taskGroupMonitor.DomainObject.Cron = taskGroupDO.Cron
		taskGroupMonitor.DomainObject.IsEnable = taskGroupDO.IsEnable

		// 通知协议，有更新
		taskGroupMonitor.Notify()
	})
}
