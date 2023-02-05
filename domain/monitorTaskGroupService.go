package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
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
	// 新的任务组不再当前列表，说明被其它节点处理了。
	if !taskGroupList.ContainsKey(taskGroupDO.Name) {
		monitor := newMonitor(taskGroupDO)
		taskGroupList.Add(taskGroupDO.Name, monitor)
		flog.Infof("任务组：%s %s 加入调度线程", taskGroupDO.Name)

		go monitor.Start()
	} else {
		// 将最新的任务组数据发送到通道
		taskGroupList.GetValue(taskGroupDO.Name).pushTaskGroup(taskGroupDO)
	}
}

// MonitorClientPush 将最新的客户端信息，推送到监控线程
func MonitorClientPush(clientDO *client.DomainObject) {
	for i := 0; i < len(clientDO.Jobs); i++ {
		if taskGroupList.ContainsKey(clientDO.Jobs[i].Name) {
			// 将最新的任务组数据发送到通道
			taskGroupList.GetValue(clientDO.Jobs[i].Name).pushClient(clientDO)
		}
	}
}

// TaskGroupMonitor 等待任务执行
type TaskGroupMonitor struct {
	TaskStatusEventBus   core.IEvent                            `inject:"TaskStatus"`   // 任务调度事件
	CheckWorkingEventBus core.IEvent                            `inject:"CheckWorking"` // 检查进行中的任务
	clients              collections.List[*client.DomainObject] // 客户端列表
	*taskGroup.DomainObject
	taskGroupChan chan *taskGroup.DomainObject
	clientChan    chan *client.DomainObject
}

// newMonitor 新建任务组监听器
func newMonitor(do *taskGroup.DomainObject) *TaskGroupMonitor {
	return container.ResolveIns(&TaskGroupMonitor{
		DomainObject:  do,
		taskGroupChan: make(chan *taskGroup.DomainObject, 1000),
		clientChan:    make(chan *client.DomainObject, 1000),
		clients:       collections.NewList[*client.DomainObject](),
	})
}

// pushTaskGroup 推送新最新信息
func (receiver *TaskGroupMonitor) pushTaskGroup(do *taskGroup.DomainObject) {
	if receiver.Name != do.Name ||
		receiver.Ver != do.Ver ||
		receiver.Caption != do.Caption ||
		receiver.StartAt != do.StartAt ||
		receiver.NextAt != do.NextAt ||
		receiver.Cron != do.Cron ||
		receiver.ActivateAt != do.ActivateAt ||
		receiver.LastRunAt != do.LastRunAt ||
		receiver.IsEnable != do.IsEnable ||
		receiver.RunSpeedAvg != do.RunSpeedAvg ||
		receiver.RunCount != do.RunCount ||
		receiver.Task.Id != do.Task.Id ||
		receiver.Task.Status != do.Task.Status {
		receiver.taskGroupChan <- do
	}
}

// pushTaskGroup 推送新最新信息
func (receiver *TaskGroupMonitor) pushClient(do *client.DomainObject) {
	receiver.clientChan <- do
}

// Start 监听任务组
func (receiver *TaskGroupMonitor) Start() {
	for {
		// 任务组状态不可用时，不需要调度
		for !receiver.IsEnable {
			receiver.updateTaskGroup(<-receiver.taskGroupChan)
		}
		// 没有可用客户端时，不需要调度
		for receiver.CanScheduleClient() == 0 {
			receiver.updateClient(<-receiver.clientChan)
		}

		select {
		case <-time.After(receiver.StartAt.Sub(time.Now())): // 时间到了，可以开始计算任务执行赶时间
			switch receiver.Task.Status {
			case enum.None, enum.ScheduleFail: // 如果调度失败状态，需要重新调度
				// 等待时间达了之后，开始调度
				receiver.waitScheduler()
			case enum.Scheduling:
				// 等待更新即可
				receiver.updateTaskGroup(<-receiver.taskGroupChan)
			case enum.Working:
				// 已成功调度到客户端，需要等待客户端上报状态
				receiver.waitWorking()
			case enum.Fail, enum.Success:
				// 等待更新
				receiver.updateTaskGroup(<-receiver.taskGroupChan)
			}
		case newData := <-receiver.taskGroupChan: // 任务组有更新
			receiver.updateTaskGroup(newData)
		case newData := <-receiver.clientChan: // 客户端有更新
			receiver.updateClient(newData)
		}
	}
}

// 等待调度
func (receiver *TaskGroupMonitor) waitScheduler() {
	select {
	case <-time.After(receiver.NextAt.Sub(time.Now())): // 时间到了，需要调度
		// 标记为调度中，阻止当前监听逻辑重复执行，否则会不停的重复执行调度
		receiver.Task.Status = enum.Scheduling
		receiver.TaskStatusEventBus.Publish(receiver)
	case newData := <-receiver.taskGroupChan: // 任务组有更新
		receiver.updateTaskGroup(newData)
	case newData := <-receiver.clientChan: // 客户端有更新
		receiver.updateClient(newData)
	}
}

// 等待完成
func (receiver *TaskGroupMonitor) waitWorking() {
	select {
	case <-time.After(10 * time.Second): // 每隔10秒，主动向客户端询问任务状态
		receiver.CheckWorkingEventBus.Publish(receiver)
	case newData := <-receiver.taskGroupChan: // 任务组有更新
		receiver.updateTaskGroup(newData)
	case newData := <-receiver.clientChan: // 客户端有更新
		receiver.updateClient(newData)
	}
}

// 有更新
func (receiver *TaskGroupMonitor) updateTaskGroup(newData *taskGroup.DomainObject) {
	receiver.DomainObject = newData
}

// 有更新
func (receiver *TaskGroupMonitor) updateClient(newData *client.DomainObject) {
	// 如果客户端离线了，则要移除
	if newData.Status == enum.Offline {
		receiver.clients.RemoveAll(func(item *client.DomainObject) bool {
			return item.Id == newData.Id
		})
		return
	}

	// 存在则直接更新，不存在则添加
	if !receiver.clients.Where(func(item *client.DomainObject) bool {
		if item.Id == newData.Id {
			item = newData
		}
		return item.Id == newData.Id
	}).Any() {
		receiver.clients.Add(newData)
	}
}

// GetClient 轮询的方式取到客户端
func (receiver *TaskGroupMonitor) GetClient() *client.DomainObject {
	// 使用轮询方式，根据调度时间排序，取最晚没调度的客户端
	return receiver.clients.Where(func(item *client.DomainObject) bool {
		return item.Status == enum.Scheduler
	}).OrderBy(func(item *client.DomainObject) any {
		return item.ScheduleAt
	}).First()
}

// CanScheduleClient 能调度的客户端
func (receiver *TaskGroupMonitor) CanScheduleClient() int {
	return receiver.clients.Where(func(item *client.DomainObject) bool {
		return item.Status == enum.Scheduler
	}).Count()
}
