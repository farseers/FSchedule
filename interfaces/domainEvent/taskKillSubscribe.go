package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"
	"encoding/json"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/trace"
)

// TaskKillSubscribe 停止任务（Redis订阅）
func TaskKillSubscribe(message any, _ core.EventArgs) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	var taskGroupDO taskGroup.DomainObject
	err := json.Unmarshal([]byte(message.(string)), &taskGroupDO)
	if err != nil {
		return
	}

	// 通知处理该任务组的服务端，需要调用客户端发起Kill请求
	if taskGroupMonitor := domain.GetTaskGroupMonitor(taskGroupDO.Name); taskGroupMonitor != nil {
		taskGroupMonitor.TaskKill()
	}
}
