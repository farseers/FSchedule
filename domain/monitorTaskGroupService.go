package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"time"
)

// 加入到监控的列表
var taskGroupList = collections.NewDictionary[string, *TaskGroupMonitor]()

// MonitorTaskGroupPush 将最新的任务组信息，推送到监控线程
func MonitorTaskGroupPush(taskGroupDO *taskGroup.DomainObject) {
	flog.Debugf("任务组更新通知：%s Ver%d", taskGroupDO.Name, taskGroupDO.Ver)
	// 新的任务组不再当前列表，说明被其它节点处理了。
	if !taskGroupList.ContainsKey(taskGroupDO.Name) {
		monitor := newMonitor(taskGroupDO)
		taskGroupList.Add(taskGroupDO.Name, monitor)
		flog.Infof("任务组：%s ver%d 加入调度线程", taskGroupDO.Name, taskGroupDO.Ver)

		go monitor.Start()
	} else {
		taskGroupMonitor := taskGroupList.GetValue(taskGroupDO.Name)
		*taskGroupMonitor.DomainObject = *taskGroupDO
		taskGroupMonitor.updated <- struct{}{}
	}
}

// ClientJoin 推送新最新信息
func ClientJoin(clientDO *client.DomainObject) {
	for i := 0; i < clientDO.Jobs.Count(); i++ {
		// 找到客户端支持的任务组
		jobName := clientDO.Jobs.Index(i).Name
		if taskGroupList.ContainsKey(jobName) {
			taskGroupMonitor := taskGroupList.GetValue(jobName)
			taskGroupMonitor.addClient(clientDO)
		}
	}
}

// ClientOffline 客户端移除
func ClientOffline(clientDO *client.DomainObject) {
	for i := 0; i < clientDO.Jobs.Count(); i++ {
		// 找到客户端支持的任务组
		if taskGroupList.ContainsKey(clientDO.Jobs.Index(i).Name) {
			taskGroupMonitor := taskGroupList.GetValue(clientDO.Jobs.Index(i).Name)
			taskGroupMonitor.removeClient(clientDO)
		}
	}
}

// TaskGroupMonitor 等待任务执行
type TaskGroupMonitor struct {
	SchedulerEventBus    core.IEvent                            `inject:"TaskScheduler"` // 任务调度事件
	FinishEventBus       core.IEvent                            `inject:"TaskFinish"`    // 任务完成
	CheckWorkingEventBus core.IEvent                            `inject:"CheckWorking"`  // 检查进行中的任务
	lock                 core.ILock                             // 锁
	clients              collections.List[*client.DomainObject] // 客户端列表
	updated              chan struct{}                          // 数据有更新，让流程重置
	*taskGroup.DomainObject
}

// newMonitor 新建任务组监听器
func newMonitor(do *taskGroup.DomainObject) *TaskGroupMonitor {
	return container.ResolveIns(&TaskGroupMonitor{
		DomainObject: do,
		updated:      make(chan struct{}, 1000),
		clients:      collections.NewList[*client.DomainObject](),
		lock:         container.Resolve[schedule.Repository]().NewLock(do.Name),
	})
}

// Start 监听任务组
func (receiver *TaskGroupMonitor) Start() {
	for {
		// 任务组状态不可用、没有可用客户端，不需要调度
		for !receiver.IsEnable || receiver.CanScheduleClient() == 0 {
			<-receiver.updated
			continue
		}

		select {
		case <-time.After(receiver.StartAt.Sub(time.Now())): // 时间到了，可以开始计算任务执行赶时间
			switch receiver.Task.Status {
			case enum.None, enum.ScheduleFail: // 如果调度失败状态，需要重新调度
				// 等待时间达了之后，开始调度
				receiver.waitScheduler()
			case enum.Scheduling:
				// 等待更新即可
				<-receiver.updated
			case enum.Working:
				// 已成功调度到客户端，需要等待客户端上报状态
				receiver.waitWorking()
			case enum.Fail, enum.Success:
				receiver.taskFinish()
			}
		case <-receiver.updated:
		}
	}
}

// 等待调度
func (receiver *TaskGroupMonitor) waitScheduler() {
	select {
	case <-time.After(receiver.NextAt.Sub(time.Now())): // 时间到了，需要调度
		// 标记为调度中，阻止当前监听逻辑重复执行，否则会不停的重复执行调度
		receiver.lock.TryLockRun(func() {
			receiver.Task.Scheduling()
			_ = receiver.SchedulerEventBus.Publish(receiver)
		})
	case <-receiver.updated:
	}
}

// 等待完成
func (receiver *TaskGroupMonitor) waitWorking() {
	select {
	case <-time.After(5 * time.Second): // 每隔5秒，主动向客户端询问任务状态
		receiver.lock.TryLockRun(func() {
			_ = receiver.CheckWorkingEventBus.Publish(receiver)
		})
	case <-receiver.updated:
	}
}

// 任务完成
func (receiver *TaskGroupMonitor) taskFinish() {
	receiver.lock.TryLockRun(func() {
		_ = receiver.SchedulerEventBus.Publish(receiver)
		// 等待更新
		<-receiver.updated
	})
}

// 有新客户端
func (receiver *TaskGroupMonitor) addClient(newData *client.DomainObject) {
	receiver.clients.Add(newData)
	receiver.updated <- struct{}{}
}

// 移除客户端
func (receiver *TaskGroupMonitor) removeClient(newData *client.DomainObject) {
	receiver.clients.RemoveAll(func(item *client.DomainObject) bool {
		return item.Id == newData.Id
	})
	receiver.updated <- struct{}{}
}

// PollingClient 轮询的方式取到客户端
func (receiver *TaskGroupMonitor) PollingClient() *client.DomainObject {
	// 使用轮询方式，根据调度时间排序，取最晚没调度的客户端
	return receiver.clients.Where(func(item *client.DomainObject) bool {
		return item.Status == enum.Scheduler && item.Jobs.Where(func(jobVO client.JobVO) bool {
			return jobVO.Name == receiver.Name && jobVO.Ver == receiver.Ver
		}).Any()
	}).OrderBy(func(item *client.DomainObject) any {
		return item.ScheduleAt
	}).First()
}

// GetClient 获取客户端
func (receiver *TaskGroupMonitor) GetClient() *client.DomainObject {
	// 使用轮询方式，根据调度时间排序，取最晚没调度的客户端
	return receiver.clients.Where(func(item *client.DomainObject) bool {
		return item.Id == receiver.Task.Client.Id
	}).First()
}

// CanScheduleClient 能调度的客户端
func (receiver *TaskGroupMonitor) CanScheduleClient() int {
	return receiver.clients.Where(func(item *client.DomainObject) bool {
		return item.Status == enum.Scheduler && item.Jobs.Where(func(jobVO client.JobVO) bool {
			return jobVO.Name == receiver.Name && jobVO.Ver == receiver.Ver
		}).Any()
	}).Count()
}

// TaskGroupCount 返回当前正在监控的任务组数量
func TaskGroupCount() int {
	return taskGroupList.Count()
}
