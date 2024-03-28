package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
)

// CheckWorkingEvent 检查进行中的任务
func CheckWorkingEvent(message any, _ core.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	taskGroupRepository := container.Resolve[taskGroup.Repository]()

	if do.Task.ExecuteStatus != executeStatus.Working {
		return
	}

	// 得到当前处理的客户端
	clientDO := do.GetClient()

	// 客户端下线了
	if clientDO == nil || clientDO.IsNil() || clientDO.IsOffline() {
		do.Task.SetFail("客户端下线了")
		taskGroupRepository.Save(*do.DomainObject)
		return
	}

	// todo 这里的dto如果为nil，不能立马认定为失败
	// 主动向客户端查询任务状态
	if dto, err := clientDO.CheckTaskStatus(do.Task.Id); err == nil {
		if dto.IsNil() {
			do.ReportFail(taskGroupRepository, "客户端dto返回nil")
		} else {
			do.Report(dto.Status, dto.Data, dto.Progress, dto.RunSpeed, dto.NextTimespan, "", taskGroupRepository)
		}
	}
}
