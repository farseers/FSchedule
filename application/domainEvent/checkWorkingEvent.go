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
	do := message.(*taskGroup.Monitor)
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	clientRepository := container.Resolve[client.Repository]()
	clientCheck := container.Resolve[client.IClientCheck]()

	// 取出任务组信息
	taskGroupDO := taskGroupRepository.ToEntity(do.Name)
	if taskGroupDO.Task.Status != enum.Working {
		return
	}
	// 获取任务信息
	taskEO := taskGroupRepository.GetTask(taskGroupDO.Task.Id)
	// 得到当前处理的客户端
	clientDO := clientRepository.ToEntity(taskGroupDO.Task.Client.Id)

	// 客户端下线了
	if clientDO.IsNil() {
		taskGroupDO.ClientOffline()
	} else {
		// 主动向客户端查询任务状态
		dto := clientCheck.Status(clientDO, taskGroupDO.Task.Id)
		// 更新任务
		taskEO.UpdateTask(dto.Status, dto.Data, dto.Progress, dto.RunSpeed)
	}

	domain.TaskReportService(taskEO, taskGroupDO, taskGroupRepository)
}
