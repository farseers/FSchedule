package job

import (
	"FSchedule/domain/serverNode"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/tasks"
)

// ServerActivateJob 更新节点活跃时间
func ServerActivateJob(context *tasks.TaskContext) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	repository := container.Resolve[serverNode.Repository]()
	serverNodeDO := repository.ToEntity(core.AppId)
	serverNodeDO.Activate()
	repository.Save(&serverNodeDO)
}
