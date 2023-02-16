package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/mapper"
	"time"
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
		if !do.CanScheduler() {
			do.ScheduleFail()
			return
		}

		// 轮询的方式取到客户端
		clientSchedule := do.PollingClient()
		// 没有可调度的客户端
		if clientSchedule == nil || clientSchedule.IsNil() {
			do.ScheduleFail()
			taskGroupRepository.Save(*do.DomainObject)
			return
		}

		// 分配客户端
		do.SetClient(mapper.Single[taskGroup.ClientVO](clientSchedule))

		// 请求客户端
		clientTask := mapper.Single[client.TaskEO](do.Task)
		if clientSchedule.Schedule(&clientTask) {
			// 调度成功
			clientRepository.Save(clientSchedule)
			taskGroupRepository.SaveAndTask(*do.DomainObject)
			return
		}
		// 调度失败
		do.ScheduleFail()
		clientRepository.Save(clientSchedule)
		taskGroupRepository.Save(*do.DomainObject)

		time.Sleep(100 * time.Millisecond)
	}
}
