package domainEvent

import (
	"FSchedule/domain/client"
	"FSchedule/domain/monitor"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/mapper"
)

// ScheduleEvent 任务调度
func ScheduleEvent(message any, _ eventBus.EventArgs) {
	do := message.(monitor.DomainObject)
	clientRepository := container.Resolve[client.Repository]()
	// 找到可调度的客户端
	clients := clientRepository.GetClients(do.Name, do.Ver)

	if clients.Count() > 0 {
		// 使用轮询方式，根据调度时间排序，取最晚没调度的客户端
		clientSchedule := clientRepository.GetClients(do.Name, do.Ver).OrderBy(func(item client.DomainObject) any {
			return item.ScheduleAt
		}).First()

		// 生成任务
		taskGroupRepository := container.Resolve[taskGroup.Repository]()
		taskGroupDO := taskGroupRepository.ToEntity(do.Name)
		taskGroupDO.CreateTask(taskGroup.ClientVO{
			Id:   clientSchedule.Id,
			Name: clientSchedule.Name,
			Ip:   clientSchedule.Ip,
			Port: clientSchedule.Port})

		// 调度
		clientTask := mapper.Single[client.TaskEO](taskGroupDO.Task)
		if clientSchedule.Schedule(&clientTask) {
			clientRepository.Save(clientSchedule)
		}
	}
}
