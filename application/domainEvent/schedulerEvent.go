package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"FSchedule/domain/monitor"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
	"time"
)

// SchedulerEvent 任务调度
func SchedulerEvent(message any, _ eventBus.EventArgs) {
	do := message.(*monitor.DomainObject)
	// 只订阅调度状态的事件
	if do.TaskStatus != enum.Scheduler {
		return
	}
	clientRepository := container.Resolve[client.Repository]()
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	scheduleRepository := container.Resolve[schedule.Repository]()
	lock := scheduleRepository.GetLock(do.Name)
	if !lock.TryLock() {
		flog.Debugf("调度任务时加锁失败，Job=%s，serverIP=%s，serverId=%d", do.Name, fs.AppIp, fs.AppId)
		return
	}
	defer lock.ReleaseLock()

	// 找到可调度的客户端
	clients := clientRepository.GetClients(do.Name, do.Ver)
	taskGroupDO := taskGroupRepository.ToEntity(do.Name)
	if !taskGroupDO.CanScheduler() {
		return
	}

	for clients.Count() > 0 {
		// 使用轮询方式，根据调度时间排序，取最晚没调度的客户端
		clientSchedule := clientRepository.GetClients(do.Name, do.Ver).OrderBy(func(item client.DomainObject) any {
			return item.ScheduleAt
		}).First()

		// 生成任务
		taskGroupDO.CreateTask(taskGroup.ClientVO{
			Id:   clientSchedule.Id,
			Name: clientSchedule.Name,
			Ip:   clientSchedule.Ip,
			Port: clientSchedule.Port})

		// 调度
		clientTask := mapper.Single[client.TaskEO](taskGroupDO.Task)
		if clientSchedule.Schedule(&clientTask) {
			clientRepository.Save(clientSchedule)
			taskGroupRepository.Save(taskGroupDO)
			domain.MonitorPush(taskGroupDO)
			return
		}

		// 移除刚才失败的客户端
		clients.Remove(clientSchedule)
	}

	taskGroupDO.ScheduleFail()
	taskGroupRepository.Save(taskGroupDO)

	// 调度失败，最理想的做法：让集群的其它节点立即处理，当前节点不处理。
	// 所以这里等待3秒，3秒是每个任务组重新刷新数据的时间。
	time.Sleep(3 * time.Second)
	domain.MonitorPush(taskGroupDO)
}
