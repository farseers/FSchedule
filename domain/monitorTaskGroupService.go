package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"bytes"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/timingWheel"
	"runtime"
	"strconv"
	"time"
)

// 加入到监控的列表
var taskGroupList = collections.NewDictionary[string, *TaskGroupMonitor]()
var tw = timingWheel.New(100*time.Millisecond, 120)

// MonitorTaskGroupPush 将最新的任务组信息，推送到监控线程
func MonitorTaskGroupPush(taskGroupDO *taskGroup.DomainObject) {
	// 新的任务组不再当前列表，说明被其它节点处理了。
	if !taskGroupList.ContainsKey(taskGroupDO.Name) {
		monitor := newMonitor(taskGroupDO)
		taskGroupList.Add(taskGroupDO.Name, monitor)
		flog.Infof("任务组：%s ver:%d 加入调度线程", taskGroupDO.Name, taskGroupDO.Ver)
		go monitor.Start()
	} else {
		taskGroupMonitor := taskGroupList.GetValue(taskGroupDO.Name)
		//oldTaskGroup := *taskGroupMonitor.DomainObject
		*taskGroupMonitor.DomainObject = *taskGroupDO
		taskGroupMonitor.updated <- struct{}{}
		//if oldTaskGroup.Task.Status != taskGroupMonitor.DomainObject.Task.Status ||
		//	oldTaskGroup.Task.Id != taskGroupMonitor.DomainObject.Task.Id ||
		//	oldTaskGroup.IsEnable != taskGroupMonitor.DomainObject.IsEnable {
		//	flog.Debugf("任务组更新通知：%s Ver:%d", taskGroupDO.Name, taskGroupDO.Ver)
		//	taskGroupMonitor.updated <- struct{}{}
		//}
	}
}

// ClientUpdate 客户端有更新，推送通知
func ClientUpdate(clientDO *client.DomainObject) {
	flog.Debugf("客户端（%d）更新通知：%s:%d", clientDO.Id, clientDO.Ip, clientDO.Port)
	for i := 0; i < clientDO.Jobs.Count(); i++ {
		// 找到客户端支持的任务组
		jobName := clientDO.Jobs.Index(i).Name
		if taskGroupList.ContainsKey(jobName) {
			taskGroupMonitor := taskGroupList.GetValue(jobName)
			taskGroupMonitor.updateClient(clientDO)
		}
	}
}

// TaskGroupMonitor 等待任务执行
type TaskGroupMonitor struct {
	SchedulerEventBus    core.IEvent                                         `inject:"TaskScheduler"` // 任务调度事件
	FinishEventBus       core.IEvent                                         `inject:"TaskFinish"`    // 任务完成
	CheckWorkingEventBus core.IEvent                                         `inject:"CheckWorking"`  // 检查进行中的任务
	lock                 core.ILock                                          // 锁
	clients              collections.Dictionary[int64, *client.DomainObject] // 客户端列表
	updated              chan struct{}                                       // 数据有更新，让流程重置
	curClient            *client.DomainObject                                // 当前调度的客户端
	*taskGroup.DomainObject
}

// newMonitor 新建任务组监听器
func newMonitor(do *taskGroup.DomainObject) *TaskGroupMonitor {
	return container.ResolveIns(&TaskGroupMonitor{
		DomainObject: do,
		updated:      make(chan struct{}, 1000),
		clients:      collections.NewDictionary[int64, *client.DomainObject](),
		lock:         container.Resolve[schedule.Repository]().NewLock(do.Name),
	})
}

// Start 监听任务组
func (receiver *TaskGroupMonitor) Start() {
	for {
		// 清空更新队列
		for len(receiver.updated) > 0 {
			<-receiver.updated
		}
		switch receiver.Task.Status {
		case enum.None, enum.ScheduleFail: // 如果调度失败状态，需要重新调度
			// 等待时间达了之后，开始调度
			receiver.waitStart()
		case enum.Scheduling:
			// 等待更新即可
			flog.Debugf("任务组：%s 等待更新", receiver.Name)
			<-receiver.updated
		case enum.Working:
			// 已成功调度到客户端，等待客户端执行完成
			receiver.waitWorking()
		case enum.Fail, enum.Success:
			receiver.taskFinish()
		}
	}
}

// 等待开始
func (receiver *TaskGroupMonitor) waitStart() {
	for {
		if receiver.Task.Status != enum.None && receiver.Task.Status != enum.ScheduleFail {
			return
		}

		// 任务组状态不可用、没有可用客户端，不需要调度
		if !receiver.IsEnable {
			flog.Debugf("任务组：%s "+flog.YellowString("停止状态，等待任务重新开启"), receiver.Name)
			<-receiver.updated
			continue
		}

		// 任务组状态不可用、没有可用客户端，不需要调度
		if receiver.CanScheduleClient() == 0 {
			flog.Debugf("任务组：%s "+flog.YellowString("没有客户端，等待客户端接入"), receiver.Name)
			<-receiver.updated
			continue
		}

		flog.Debugf("任务组：%s 等待开始时间", receiver.Name)
		timer := tw.AddTimePrecision(receiver.StartAt)
		select {
		case <-timer.C: // 开始时间到了，可以开始计算任务执行赶时间
			flog.Debugf("任务组：%s 等待执行时间", receiver.Name)
			receiver.waitScheduler()
			return
		case <-receiver.updated:
			timer.Stop()
		}
	}
}

// 等待调度
func (receiver *TaskGroupMonitor) waitScheduler() {
	// 由于创建锁的时候，需要网络IO开销，所以这里提前50ms进入
	select {
	case <-tw.AddTime(receiver.Task.StartAt.Add(-50 * time.Millisecond)).C: // 执行时间到了，准开始调度
		flog.Debugf("任务组：%s %d 准备抢占锁，延迟：%d us GoID:%d", receiver.Name, receiver.Task.Id, time.Since(receiver.Task.StartAt).Microseconds(), GoID())
		// 标记为调度中，阻止当前监听逻辑重复执行，否则会不停的重复执行调度
		if !receiver.lock.TryLockRun(func() {
			flog.Debugf("任务组：%s %d 抢锁成功，延迟：%d us GoID:%d", receiver.Name, receiver.Task.Id, time.Since(receiver.Task.StartAt).Microseconds(), GoID())
			// 前面，提前了50ms进到这里。所以这里需要转成正确的时间
			<-tw.AddTimePrecision(receiver.Task.StartAt).C
			receiver.Task.Scheduling()
			flog.Debugf("任务组：%s %d 发布调度事件，延迟：%d us GoID:%d", receiver.Name, receiver.Task.Id, time.Since(receiver.Task.StartAt).Microseconds(), GoID())
			_ = receiver.SchedulerEventBus.Publish(receiver)
		}) {
			flog.Debugf("任务组：%s %d 没有抢到锁，延迟：%d us GoID:%d", receiver.Name, receiver.Task.Id, time.Since(receiver.Task.StartAt).Microseconds(), GoID())
			// 没有抢到锁，就等更新
			<-receiver.updated
		}
	case <-receiver.updated:
		flog.Debugf("任务组：%s %d 有更新", receiver.Name, receiver.Task.Id)
	}
}
func GoID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// 等待完成
func (receiver *TaskGroupMonitor) waitWorking() {
	if receiver.curClient == nil || receiver.curClient.IsNil() || receiver.curClient.IsOffline() {
		flog.Debugf("任务组：%s 当前客户端已离线", receiver.Name)
		receiver.lock.TryLockRun(func() {
			_ = receiver.CheckWorkingEventBus.Publish(receiver)
		})
	}

	flog.Debugf("任务组：%s 等待客户端执行完成", receiver.Name)
	timer := tw.Add(60 * time.Second)
	// 这里用循环是为了，任何的更新，如果仍处于Working状态，则不需要跳到外面重新执行
	for {
		select {
		case <-timer.C: // 每隔60秒，主动向客户端询问任务状态
			flog.Debugf("任务组：%s 主动向客户端询问任务状态", receiver.Name)
			receiver.lock.TryLockRun(func() {
				_ = receiver.CheckWorkingEventBus.Publish(receiver)
			})
		case <-receiver.updated:
			if !receiver.Task.IsWorking() {
				timer.Stop()
				return
			}
		}
	}
}

// 任务完成
func (receiver *TaskGroupMonitor) taskFinish() {
	flog.Debugf("任务组：%s 任务完成", receiver.Name)
	if !receiver.lock.TryLockRun(func() {
		_ = receiver.FinishEventBus.Publish(receiver.DomainObject)
	}) {
		// 没有抢到锁，就等更新
		<-receiver.updated
	}
}

// 更新客户端
func (receiver *TaskGroupMonitor) updateClient(newData *client.DomainObject) {
	flog.Debugf("任务组：%s 更新客户端updateClient", receiver.Name)
	// 状态为不可调度时，则移除列表
	if newData.IsNotSchedule() {
		receiver.clients.Remove(newData.Id)
	} else {
		receiver.clients.Add(newData.Id, newData)
		receiver.updated <- struct{}{}
	}
}

// PollingClient 轮询的方式取到客户端
func (receiver *TaskGroupMonitor) PollingClient() *client.DomainObject {
	lst := receiver.clients.Values()
	for ver := receiver.Ver; ver > 0; ver-- {
		// 使用轮询方式，根据调度时间排序，取最晚没调度的客户端
		receiver.curClient = lst.Where(func(item *client.DomainObject) bool {
			return item.Status == enum.Scheduler && item.Jobs.Where(func(jobVO client.JobVO) bool {
				return jobVO.Name == receiver.Name && jobVO.Ver == ver
			}).Any()
		}).OrderBy(func(item *client.DomainObject) any {
			return item.ScheduleAt.UnixMilli()
		}).First()

		// 找到了，不用继续往下找
		if receiver.curClient != nil {
			break
		}
	}
	return receiver.curClient
}

// GetClient 获取客户端
func (receiver *TaskGroupMonitor) GetClient() *client.DomainObject {
	return receiver.curClient
}

// CanScheduleClient 能调度的客户端
func (receiver *TaskGroupMonitor) CanScheduleClient() int {
	return receiver.clients.Count()
}

// TaskGroupCount 返回当前正在监控的任务组数量
func TaskGroupCount() int {
	for _, v := range taskGroupList.ToMap() {
		flog.Debugf("任务组：%s，状态：%s", v.Name, v.Task.Status.String())
	}
	return taskGroupList.Count()
}

// TaskGroupEnableCount 返回开启状态的任务组
func TaskGroupEnableCount() int {
	return taskGroupList.Values().Where(func(item *TaskGroupMonitor) bool {
		return item.CanScheduleClient() > 0
	}).Count()
}
