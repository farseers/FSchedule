package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
)

// CheckWorkingEvent 检查进行中的任务
func CheckWorkingEvent(message any, _ core.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	clientRepository := container.Resolve[client.Repository]()
	clientCheck := container.Resolve[client.IClientCheck]()

	if do.Task.Status != enum.Working {
		return
	}

	// 得到当前处理的客户端
	clientDO := do.GetClient()

	// 客户端下线了
	if clientDO == nil || clientDO.IsNil() || clientDO.IsOffline() {
		do.ClientOffline()
		taskGroupRepository.Save(*do.DomainObject)
		return
	}

	// 主动向客户端查询任务状态
	dto, err := clientCheck.Status(clientDO, do.Task.Id)
	if err != nil {
		clientDO.UnSchedule()
		clientRepository.Save(clientDO)
	} else {
		if dto.Id == 0 || dto.Name == "" || dto.TaskGroupId == 0 {
			do.ReportFail(taskGroupRepository)
		} else {
			do.Report(dto.Status, dto.Data, dto.Progress, dto.RunSpeed, dto.NextTimespan, taskGroupRepository)
		}
	}
}
