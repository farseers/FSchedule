package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum/clientStatus"
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"context"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/timingWheel"
	"time"
)

// 加入到监控的列表
var taskGroupList = collections.NewDictionary[string, *TaskGroupMonitor]()

// MonitorTaskGroupPush 将最新的任务组信息，推送到监控线程
func MonitorTaskGroupPush(taskGroupDO *taskGroup.DomainObject) {
	// 新的任务组不再当前列表，说明被其它节点处理了。
	if !taskGroupList.ContainsKey(taskGroupDO.Name) {
		if !taskGroupDO.IsEnable {
			return
		}
		// 加入到任务组监控列表
		monitor := newMonitor(taskGroupDO)
		taskGroupList.Add(taskGroupDO.Name, monitor)

		// 找到当前任务组支持的客户端列表，主动通知，更快速接入
		// 用在将任务组.IsEnable由false改成true时
		clientList.Values().Foreach(func(clientMonitor **ClientMonitor) {
			(*clientMonitor).client.Jobs.Foreach(func(job *client.JobVO) {
				if job.Name == taskGroupDO.Name {
					monitor.updateClient((*clientMonitor).client)
				}
			})
		})
		// 开启协程
		go monitor.Start()
	} else {
		taskGroupMonitor := taskGroupList.GetValue(taskGroupDO.Name)
		// 之前是运行状态，改为停止状态，则需要退出调度线程
		needKill := taskGroupMonitor.IsEnable && !taskGroupDO.IsEnable
		*taskGroupMonitor.DomainObject = *taskGroupDO
		taskGroupMonitor.updateNotice()
		if needKill {
			// 强制退出线程
			taskGroupMonitor.cancelFunc()
		}
	}
}

// ClientUpdate 客户端有更新，推送通知
func ClientUpdate(clientDO *client.DomainObject) {
	//flog.Debugf("客户端（%d）更新通知：%s:%d", clientDO.Id, clientDO.Ip, clientDO.Port)
	for i := 0; i < clientDO.Jobs.Count(); i++ {
		// 找到客户端支持的任务组
		jobName := clientDO.Jobs.Index(i).Name
		for _, taskGroupMonitor := range taskGroupList.ToMap() {
			if taskGroupMonitor.Name == jobName {
				taskGroupMonitor.updateClient(clientDO)
			}
		}
	}
}

// TaskGroupMonitor 等待任务执行
type TaskGroupMonitor struct {
	SchedulerEventBus    core.IEvent                                         `inject:"TaskScheduler"` // 任务调度事件
	FinishEventBus       core.IEvent                                         `inject:"TaskFinish"`    // 任务完成
	CheckWorkingEventBus core.IEvent                                         `inject:"CheckWorking"`  // 检查进行中的任务
	ScheduleRepository   schedule.Repository                                 // 锁
	clients              collections.Dictionary[int64, *client.DomainObject] // 客户端列表
	updated              chan struct{}                                       // 数据有更新，让流程重置
	curClient            *client.DomainObject                                // 当前调度的客户端
	isWorking            bool                                                // 是否进入工作状态
	waitWork             bool                                                // 进入抢占锁状态
	*taskGroup.DomainObject
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// newMonitor 新建任务组监听器
func newMonitor(do *taskGroup.DomainObject) *TaskGroupMonitor {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return container.ResolveIns(&TaskGroupMonitor{
		DomainObject: do,
		updated:      make(chan struct{}, 1000),
		clients:      collections.NewDictionary[int64, *client.DomainObject](),
		ctx:          ctx,
		cancelFunc:   cancelFunc,
	})
}

// Start 监听任务组
func (receiver *TaskGroupMonitor) Start() {
	// 退出时，移除监控
	defer func() {
		receiver.isWorking = false
		receiver.waitWork = false
		taskGroupList.Remove(receiver.Name)
		flog.Infof("任务组：%s ver:%s 退出调度线程", flog.Blue(receiver.Name), flog.Yellow(receiver.Ver))
	}()

	// 抢占锁，谁抢到，谁负责这个任务组的调度（只允许一个集群节点监控任务组）
	receiver.waitWork = true
	receiver.ScheduleRepository.Schedule(receiver.Name, func() {
		receiver.isWorking = true
		flog.Infof("任务组：%s ver:%s 加入调度线程", flog.Blue(receiver.Name), flog.Yellow(receiver.Ver))
		for {
			// 清空更新队列
			receiver.updated = make(chan struct{}, 1000)

			select {
			case <-receiver.ctx.Done(): // 任务组停止，或删除时退出
				return
			default: // 没有停止时，继续往下走
				if !receiver.IsEnable { // // 当corn格式错误时，会强制设为false，等待手动启动
					flog.Infof("任务组：%s ver:%s 状态未启用", flog.Blue(receiver.Name), flog.Yellow(receiver.Ver))
					return
				}
			}

			switch receiver.Task.ScheduleStatus {
			// 如果调度失败状态，需要重新调度
			case scheduleStatus.None:
				receiver.waitStart()
			case scheduleStatus.Scheduling:
				// 等待其它协程更新状态
				<-receiver.updated
			case scheduleStatus.Fail:
				receiver.taskFinish()
			case scheduleStatus.Success:
				switch receiver.Task.ExecuteStatus {
				case executeStatus.None:
					// 等待客户端上报运行状态
					receiver.waitJobReportWorkStatus()
				case executeStatus.Working:
					// 已成功调度到客户端，等待客户端执行完成
					receiver.waitWorking()
				case executeStatus.Fail, executeStatus.Success:
					receiver.taskFinish()
				default:
					flog.Warningf("任务组：%s ver:%s 出现未知执行状态：%d 将强制设为失败状态", flog.Blue(receiver.Name), flog.Yellow(receiver.Ver), receiver.Task.ExecuteStatus)
					receiver.Task.SetFail(fmt.Sprintf("出现未知执行状态：%d", receiver.Task.ExecuteStatus))
					receiver.taskFinish()
				}
			}
		}
	})
}

// 等待开始
func (receiver *TaskGroupMonitor) waitStart() {
	// 手动提前kill的任务，调度状态 = None，执行状态 = fail
	if receiver.Task.ExecuteStatus.IsFinish() {
		receiver.taskFinish()
		return
	}

	// 没有可用客户端，不需要调度
	if receiver.CanScheduleClient() == 0 {
		select {
		case <-receiver.updated:
			// 有可能enable有变，所以这里要返回出去，让外面来判断
			return
		}
	}

	timer := timingWheel.AddTimePrecision(receiver.StartAt.ToTime())
	select {
	case <-timer.C: // 开始时间到了，可以开始计算任务执行赶时间
		receiver.waitScheduler()
	case <-receiver.updated:
		timer.Stop()
	}
}

// 等待调度
func (receiver *TaskGroupMonitor) waitScheduler() {
	// 由于创建锁的时候，需要网络IO开销，所以这里提前100ms进入
	timer := timingWheel.AddTime(receiver.Task.StartAt.AddMillisecond(-100).ToTime())
	select {
	case <-timer.C: // 执行时间到了，准开始调度
		// 提前了100ms进到这里。
		receiver.Task.SetScheduling()
		// 发布调度事件
		_ = receiver.SchedulerEventBus.Publish(receiver)
	case <-receiver.updated:
		timer.Stop()
	}
}

// 等待完成
func (receiver *TaskGroupMonitor) waitWorking() {
	if receiver.curClient == nil || receiver.curClient.IsNil() || receiver.curClient.IsOffline() {
		flog.Debugf("任务组：%s 当前客户端已离线", receiver.Name)
		_ = receiver.CheckWorkingEventBus.Publish(receiver)
		return
	}

	// 小于10s，则按10s检查一次
	checkWorkWaitMillisecond := time.Duration(receiver.RunSpeedAvg * 3)
	if checkWorkWaitMillisecond < 10000 {
		checkWorkWaitMillisecond = 10000
	}
	timer := timingWheel.Add(checkWorkWaitMillisecond * time.Millisecond)
	select {
	case <-timer.C: // 每隔60秒，主动向客户端询问任务状态
		flog.Debugf("任务组：%s 主动向客户端询问任务状态", receiver.Name)
		_ = receiver.CheckWorkingEventBus.Publish(receiver)
	case <-receiver.updated:
		timer.Stop()
	}
}

// 任务完成
func (receiver *TaskGroupMonitor) taskFinish() {
	_ = receiver.FinishEventBus.Publish(receiver.DomainObject)
}

// 检查超时未执行
func (receiver *TaskGroupMonitor) waitJobReportWorkStatus() {
	timer := timingWheel.AddTime(receiver.Task.SchedulerAt.AddSeconds(10).ToTime())
	select {
	case <-timer.C: // 调度时间超过10s，仍未执行
		_ = receiver.CheckWorkingEventBus.Publish(receiver)
	case <-receiver.updated:
		timer.Stop()
	}
}

// 更新客户端
func (receiver *TaskGroupMonitor) updateClient(newData *client.DomainObject) {
	// 状态为不可调度时，则移除列表
	if newData.IsNotSchedule() {
		// 移除客户端
		if receiver.clients.ContainsKey(newData.Id) {
			receiver.clients.Remove(newData.Id)

			if receiver.curClient != nil && receiver.curClient.Id == newData.Id {
				receiver.curClient = nil
			}

			receiver.updateNotice()
			flog.Debugf("任务组：%s 移除客户端：%s %d 状态：%s", receiver.Name, newData.Name, newData.Id, newData.Status.String())
		}
	} else {
		if !receiver.clients.ContainsKey(newData.Id) {
			receiver.clients.Add(newData.Id, newData)
			receiver.updateNotice()
			flog.Debugf("任务组：%s 添加客户端：%s %d", receiver.Name, newData.Name, newData.Id)
		}
	}
}

// PollingClient 轮询的方式取到客户端
func (receiver *TaskGroupMonitor) PollingClient() *client.DomainObject {
	lst := receiver.clients.Values()
	for ver := receiver.Ver; ver > 0; ver-- {
		// 使用轮询方式，根据调度时间排序，取最晚没调度的客户端
		receiver.curClient = lst.Where(func(item *client.DomainObject) bool {
			return item.Status == clientStatus.Scheduler && item.Jobs.Where(func(jobVO client.JobVO) bool {
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

// 通知客户端有更新
func (receiver *TaskGroupMonitor) updateNotice() {
	// 进入抢占锁状态且没有在工作，说明没有拿到执行权，不需要更新
	if !receiver.waitWork || receiver.isWorking {
		//flog.Debugf("任务组更新通知：%s Ver:%d", taskGroupDO.Name, taskGroupDO.Ver)
		receiver.updated <- struct{}{}
	}
}

// TaskGroupCount 返回当前正在监控的任务组数量
func TaskGroupCount() int {
	lstLog := collections.NewList[string]()
	for _, v := range taskGroupList.ToMap() {
		if v.clients.Count() > 0 {
			var curClientId int64
			if v.curClient != nil {
				curClientId = v.curClient.Id
			}
			lstLog.Add(fmt.Sprintf("任务组：%s，\t状态：%s，客户端%s个，当前客户端：%s", flog.Blue(v.Name), v.Task.ExecuteStatus.String(), flog.Red(v.clients.Count()), flog.Green(curClientId)))
		}
	}
	if lstLog.Count() > 0 {
		fmt.Sprintln(lstLog.ToString("\n"))
	}
	return taskGroupList.Count()
}

// TaskGroupEnableCount 返回开启状态的任务组
func TaskGroupEnableCount() int {
	return taskGroupList.Values().Where(func(item *TaskGroupMonitor) bool {
		return item.CanScheduleClient() > 0
	}).Count()
}

// 获取任务组接受调度的客户端列表
func GetClientList(taskGroupName string) collections.List[*client.DomainObject] {
	if taskGroupList.ContainsKey(taskGroupName) {
		return taskGroupList.GetValue(taskGroupName).clients.Values()
	}
	return collections.NewList[*client.DomainObject]()
}
