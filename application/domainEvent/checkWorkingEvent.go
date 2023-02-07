package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/container"
)

// CheckWorkingEvent 检查进行中的任务
func CheckWorkingEvent(message any, _ eventBus.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	clientRepository := container.Resolve[client.Repository]()
	clientCheck := container.Resolve[client.IClientCheck]()

	// 取出任务组信息
	taskGroupDO := taskGroupRepository.ToEntity(do.Name)
	if taskGroupDO.Task.Status != enum.Working {
		return
	}
	// 得到当前处理的客户端
	clientDO := clientRepository.ToEntity(taskGroupDO.Task.Client.Id)

	// 客户端下线了
	if clientDO.IsNil() {
		taskGroupDO.ClientOffline()
		taskGroupRepository.Save(taskGroupDO)
		return
	}

	// 客户端下线了
	if clientDO.IsOffline() {
		taskGroupDO.ClientOffline()
		taskGroupRepository.Save(taskGroupDO)
		return
	}

	// 主动向客户端查询任务状态
	dto, err := clientCheck.Status(clientDO, taskGroupDO.Task.Id)
	if err != nil {
		clientDO.UnSchedule()
		clientRepository.Save(clientDO)
	} else {
		taskGroupDO.Task.UpdateTask(dto.Status, dto.Data, dto.Progress, dto.RunSpeed)
		taskGroupRepository.Save(taskGroupDO)
	}
}
