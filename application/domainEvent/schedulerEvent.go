package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/mapper"
)

// SchedulerEvent 任务调度
func SchedulerEvent(message any, _ core.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	// 只订阅调度状态的事件
	if do.Task.Status != enum.Scheduling {
		return
	}
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	clientRepository := container.Resolve[client.Repository]()

	for {
		// 取出最新的任务组
		taskGroupDO := taskGroupRepository.ToEntity(do.Name)
		if !taskGroupDO.CanScheduler() {
			return
		}

		// 轮询的方式取到客户端
		clientSchedule := do.GetClient()
		// 没有可调度的客户端
		if clientSchedule.IsNil() {
			taskGroupDO.Task.ScheduleFail()
			taskGroupRepository.Save(taskGroupDO)
			return
		}

		// 分配客户端
		taskGroupDO.SetClient(mapper.Single[taskGroup.ClientVO](clientSchedule))

		// 请求客户端
		clientTask := mapper.Single[client.TaskEO](taskGroupDO.Task)
		if clientSchedule.Schedule(&clientTask) {
			// 调度成功
			clientRepository.Save(clientSchedule)
			taskGroupRepository.SaveAndTask(taskGroupDO)
			return
		} else {
			// 调度失败
			taskGroupDO.Task.ScheduleFail()
			clientRepository.Save(clientSchedule)
			taskGroupRepository.Save(taskGroupDO)
		}
	}
}
