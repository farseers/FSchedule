package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
)

// CheckWorkingEvent 检查进行中的任务
func CheckWorkingEvent(message any, _ core.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	taskGroupRepository := container.Resolve[taskGroup.Repository]()

	if do.Task.ExecuteStatus.IsFinish() {
		return
	}
	if do.Name != do.Task.Name {
		_ = flog.Errorf("任务组：%s 注意，检查进行中的任务时，发现task.Name不一致，TaskId=%d，taskName=%s, task=%+v", do.Name, do.Task.Id, do.Task.Name, do.Task)
	}
	// 得到当前处理的客户端
	clientDO := do.GetClient()

	// 客户端下线了
	if clientDO == nil || clientDO.IsNil() || clientDO.IsOffline() {
		do.ReportFail(taskGroupRepository, "客户端下线了")
		return
	}

	// 主动向客户端查询任务状态
	if dto, err := clientDO.CheckTaskStatus(do.Name, do.Task.Id); err == nil {
		if dto.IsNil() {
			do.ReportFail(taskGroupRepository, "客户端dto返回nil")
		} else {
			if do.Name != dto.Name {
				_ = flog.Errorf("任务组：%s 注意，检查进行中的任务时，发现task.Name不一致，TaskId=%d，taskName=%s, task=%+v", do.Name, do.Task.Id, do.Task.Name, do.Task)
			}
			do.Report(dto.Status, dto.Data, dto.Progress, dto.NextTimespan, "", taskGroupRepository)
		}
	}

	// 更新客户端
	container.Resolve[client.Repository]().Save(clientDO)
}
