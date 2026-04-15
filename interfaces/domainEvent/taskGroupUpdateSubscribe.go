package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"

	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/fs/trace"
)

// TaskGroupUpdateSubscribe 任务组有更新（Redis订阅）
func TaskGroupUpdateSubscribe(message any, _ core.EventArgs) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	var taskGroupDO taskGroup.DomainObject
	err := snc.Unmarshal([]byte(message.(string)), &taskGroupDO)
	if err != nil {
		flog.Warningf("收到更新请求,任务组: %s 但出错了: %s %s", taskGroupDO.Name, err.Error(), message.(string))
		return
	}

	// 通知处理该任务组的服务端，需要调用客户端发起Kill请求
	lstTaskGroupMonitor := domain.GetTaskGroupMonitorByName(taskGroupDO.Name)
	if lstTaskGroupMonitor.Count() == 0 {
		flog.Infof("收到更新请求,任务组: %s ,但当前节点没有任务组的监控列表", taskGroupDO.Name)
		return
	}

	flog.Infof("收到更新请求,任务组: %s 共%d个客户端", taskGroupDO.Name, lstTaskGroupMonitor.Count())

	for _, taskGroupMonitor := range lstTaskGroupMonitor.ToArray() {
		client := "空"
		if taskGroupMonitor.Client != nil {
			client = taskGroupMonitor.Client.Id
		}
		flog.Infof("收到更新请求,任务组: %s %s", taskGroupMonitor.Name, client)

		flog.Infof("do = %v", taskGroupMonitor.DomainObject == nil)
		// 之前是运行状态，改为停止状态，则需要退出调度线程
		if taskGroupMonitor.IsEnable && !taskGroupDO.IsEnable {
			// 主动通知客户端，停止任务
			taskGroupMonitor.TaskKill()
		}
		flog.Infof("1")
		taskGroupMonitor.DomainObject.Data = taskGroupDO.Data
		taskGroupMonitor.DomainObject.Caption = taskGroupDO.Caption
		taskGroupMonitor.DomainObject.StartAt = taskGroupDO.StartAt

		flog.Infof("2")
		// 手动修改了执行时间
		if taskGroupMonitor.DomainObject.NextAt.ToString("yyyy-MM-dd HH:mm:ss") != taskGroupDO.NextAt.ToString("yyyy-MM-dd HH:mm:ss") {
			taskGroupMonitor.DomainObject.NextAt = taskGroupDO.NextAt
			taskGroupMonitor.DomainObject.Task.StartAt = taskGroupDO.NextAt
		}

		flog.Infof("3")
		taskGroupMonitor.DomainObject.Cron = taskGroupDO.Cron
		taskGroupMonitor.DomainObject.IsEnable = taskGroupDO.IsEnable

		flog.Infof("4")
		// 通知协议，有更新
		taskGroupMonitor.Notify()
	}
}
