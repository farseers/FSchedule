package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
)

// SchedulerEvent 任务调度
func SchedulerEvent(message any, _ eventBus.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	// 只订阅调度状态的事件
	if do.Task.Status != enum.Scheduling {
		return
	}
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	scheduleRepository := container.Resolve[schedule.Repository]()
	clientRepository := container.Resolve[client.Repository]()
	lock := scheduleRepository.GetLock(do.Name)
	if !lock.TryLock() {
		flog.Debugf("调度任务时加锁失败，Job=%s，serverIP=%s，serverId=%d", do.Name, fs.AppIp, fs.AppId)
		return
	}
	defer lock.ReleaseLock()

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
