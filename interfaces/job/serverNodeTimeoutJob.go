package job

import (
	"FSchedule/domain/serverNode"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/tasks"
)

// ServerNodeTimeoutJob 移除30秒不活跃的
func ServerNodeTimeoutJob(context *tasks.TaskContext) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	repository := container.Resolve[serverNode.Repository]()
	repository.ToList().Foreach(func(serverNodeDO *serverNode.DomainObject) {
		if dateTime.Since(serverNodeDO.ActivateAt).Seconds() >= 30 {
			repository.Remove(serverNodeDO.Id)
			flog.Infof("集群节点：%s %s:%d 不再活跃，移出集群", flog.Green(serverNodeDO.Id), flog.Yellow(serverNodeDO.Ip), serverNodeDO.Port)
		}
	})
}
