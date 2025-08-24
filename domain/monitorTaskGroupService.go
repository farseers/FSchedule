package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"fmt"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/color"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/timingWheel"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
)

// 加入到监控的列表
var taskGroupList = collections.NewDictionary[string, *TaskGroupMonitor]()

// 找到该任务组的监控
func GetTaskGroupMonitor(clientId string) *TaskGroupMonitor {
	return taskGroupList.GetValue(clientId)
}

// 找到该任务组的监控
func GetTaskGroupMonitorByName(taskGroupName string) collections.List[*TaskGroupMonitor] {
	lst := taskGroupList.Values()
	lst.RemoveAll(func(item *TaskGroupMonitor) bool {
		return item.Name != taskGroupName
	})
	return lst
}

// 移除单个客户端任务组监控
func RemoveMonitorClient(clientId string) {
	taskGroupMonitor := GetTaskGroupMonitor(clientId)
	taskGroupList.Remove(clientId)
	if taskGroupMonitor != nil {
		if taskGroupMonitor.Client != nil {
			taskGroupMonitor.Client.Close()
			taskGroupMonitor.Client.IsMaster = false
			container.Resolve[client.Repository]().Save(*taskGroupMonitor.Client)
			taskGroupMonitor.Client = nil
		}
	}
}

// 移除任务组监控
func RemoveMonitorTaskGroupName(taskGroupName string) {
	GetTaskGroupMonitorByName(taskGroupName).Foreach(func(item **TaskGroupMonitor) {
		// 找到所有客户端，然后删除
		RemoveMonitorClient((*item).Client.Id)
	})
}

// TaskGroupMonitor 等待任务执行
type TaskGroupMonitor struct {
	*taskGroup.DomainObject
	ScheduleRepository schedule.Repository  // 锁
	Client             *client.DomainObject // 客户端
	updated            chan struct{}        // 数据有更新，让流程重置
}

// MonitorTaskGroupPush 将最新的任务组信息，推送到监控线程
func MonitorTaskGroupPush(clientDO *client.DomainObject, taskGroupDO *taskGroup.DomainObject) {
	taskGroupMonitor := GetTaskGroupMonitor(clientDO.Id)

	// 新接入的任务组
	if taskGroupMonitor == nil {
		// 加入到任务组监控列表
		taskGroupMonitor = container.ResolveIns(&TaskGroupMonitor{
			DomainObject: taskGroupDO,
			updated:      make(chan struct{}, 1000),
			Client:       clientDO,
		})
		taskGroupList.Add(clientDO.Id, taskGroupMonitor)

		// 开启协程
		go taskGroupMonitor.Start()
	} else {
		// 之前是运行状态，改为停止状态，则需要退出调度线程
		needKill := taskGroupMonitor.IsEnable && !taskGroupDO.IsEnable
		*taskGroupMonitor.DomainObject = *taskGroupDO
		taskGroupMonitor.updated <- struct{}{}
		if needKill {
			// 主动通知客户端，停止任务
			taskGroupMonitor.TaskKill()
		}
	}
}

// Start 监听任务组
func (receiver *TaskGroupMonitor) Start() {
	// 抢占锁，谁抢到，谁负责这个任务组的调度（只允许一个集群节点监控任务组）
	receiver.ScheduleRepository.Schedule(receiver.Name, func() {
		taskGroupRepository := container.Resolve[taskGroup.Repository]()

		// 退出时，移除监控
		defer func() {
			// 如果任务组的状态是进行中，则要强制失败
			if receiver.Task.ScheduleStatus != scheduleStatus.None && !receiver.Task.IsFinish() {
				receiver.ReportFail("客户端下线了", taskGroupRepository)
				receiver.taskFinish()
			}
			if receiver.Client == nil {
				flog.Errorf("任务组：%s ver:%s 退出调度线程时 client = nil", color.Blue(receiver.Name), color.Yellow(receiver.Ver))
			} else {
				flog.Infof("任务组：%s ver:%s 客户端：%s 退出调度线程", color.Blue(receiver.Name), color.Yellow(receiver.Ver), receiver.Client.Id)
				RemoveMonitorClient(receiver.Client.Id)
			}
		}()

		// 有可能原节点挂了，由另外节点继续接管，所以需要重新取到最新的对象（因为现在取消了任务组数据的实时订阅发送）
		if receiver.DomainObject == nil { // 这里有可能为nil
			return
		}
		*receiver.DomainObject = taskGroupRepository.ToEntity(receiver.Name)
		receiver.Client.IsMaster = true
		container.Resolve[client.Repository]().Save(*receiver.Client)

		// 重新连接进来时，有可能上一次的任务执行了一半。因此这里要做检查
		if receiver.Task.ScheduleStatus != scheduleStatus.None {
			receiver.Task.SetFail("客户端重连，强制取消上次未执行的任务")
			receiver.taskFinish()
		}

		flog.Infof("任务组：%s ver:%s 加入调度线程", color.Blue(receiver.Name), color.Yellow(receiver.Ver))
		for {
			// 清空更新队列
			receiver.updated = make(chan struct{}, 1000)
			receiver.ActivateAt = dateTime.Now()

			select {
			case <-receiver.Client.Ctx.Done(): // 任务组停止，或删除时退出
				return
			default:
				// 如果任务是停止状态，则等待fops开启后继续执行
				if !receiver.IsEnable {
					// 如果任务是完成状态，则重新初始化
					if receiver.Task.IsFinish() {
						receiver.taskFinish()
					}
					<-receiver.updated
					continue
				}
			}

			switch receiver.Task.ScheduleStatus {
			// 如果调度失败状态，需要重新调度
			case scheduleStatus.None:
				receiver.waitStart()
			case scheduleStatus.Scheduling:
				timer := timingWheel.Add(5 * time.Second)
				select {
				// 5秒没反应，则认为调度超时
				case <-timer.C:
					timer.Stop()
					flog.Warningf("任务组：%s ver:%s 在等待调度时，客户端5秒内没反应，强制将任务标记为调度超时", color.Blue(receiver.Name), color.Yellow(receiver.Ver))
					receiver.Task.ScheduleFail("调度超时")
					receiver.taskFinish()
				// 等待其它协程更新状态
				case <-receiver.updated:
				}
			case scheduleStatus.Fail:
				receiver.taskFinish()
			case scheduleStatus.Success:
				switch receiver.Task.ExecuteStatus {
				case executeStatus.None, executeStatus.Working:
					select {
					// 任务组停止，或删除时退出
					case <-receiver.Client.Ctx.Done():
						flog.Warningf("任务组：%s ver:%s 在执行任务时，客户端断开连接，强制将任务标记为失败", color.Blue(receiver.Name), color.Yellow(receiver.Ver))
						receiver.Task.SetFail("客户端断开连接")
						receiver.taskFinish()
						return
					// 等待客户端上报运行状态
					case <-receiver.updated:
					}
				case executeStatus.Fail, executeStatus.Success:
					receiver.taskFinish()
				default:
					flog.Warningf("任务组：%s ver:%s 出现未知执行状态：%d 将强制设为失败状态", color.Blue(receiver.Name), color.Yellow(receiver.Ver), receiver.Task.ExecuteStatus)
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

	// 任务组总的有效时间
	timer := timingWheel.AddTimePrecision(receiver.StartAt.ToTime())
	select {
	// 任务组停止，或删除时退出
	case <-receiver.Client.Ctx.Done():
		return
	case <-receiver.updated:
		timer.Stop()
	// 开始时间到了，可以开始计算任务执行赶时间
	case <-timer.C:
		receiver.waitScheduler()
	// 如果receiver.StartAt时间过长，将导致时间计算不精准，这里做一个保护，相当于x分钟后，重新计算等待时间
	case <-time.After(3 * time.Minute):
	}
}

// 等待调度
func (receiver *TaskGroupMonitor) waitScheduler() {
	// 由于创建锁的时候，需要网络IO开销，所以这里提前100ms进入
	timer := timingWheel.AddTime(receiver.Task.StartAt.AddMillisecond(-500).ToTime())
	select {
	// 任务组停止，或删除时退出
	case <-receiver.Client.Ctx.Done():
		timer.Stop()
		return
	case <-receiver.updated:
		timer.Stop()
	case <-timer.C:
		timer.Stop()
		// 提前了100ms进到这里。
		receiver.Task.SetScheduling()
		// 调度
		receiver.schedulerEvent()
	// 如果receiver.Task.StartAt时间过长，将导致时间计算不精准，这里做一个保护，相当于x分钟后，重新计算等待时间
	case <-time.After(3 * time.Minute):
		timer.Stop()
	}
}

// SchedulerEvent 任务调度
func (receiver *TaskGroupMonitor) schedulerEvent() {
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	clientRepository := container.Resolve[client.Repository]()

	if !receiver.CanScheduler() {
		flog.Debugf("任务组：%s 条件不满足无法调度", receiver.Name)
		receiver.Task.ScheduleFail("条件不满足无法调度")
		return
	}

	// 没有可调度的客户端
	if receiver.Client == nil || receiver.Client.IsClose() {
		flog.Debugf("任务组：%s 客户端已断开连接，无法调度", receiver.Name)
		receiver.Task.ScheduleFail("客户端已断开连接，无法调度")
		return
	}

	// 请求客户端
	var err error
	clientTask := mapper.Single[client.TaskEO](receiver.Task)
	if err = receiver.Client.TrySchedule(clientTask); err == nil {
		// 调度成功，分配客户端
		receiver.Task.ScheduleSuccess(mapper.Single[taskGroup.ClientVO](receiver.Client))
		_ = container.Resolve[redis.IClient]("default").Transaction(func() {
			taskGroupRepository.SaveAndTask(*receiver.DomainObject)
			clientRepository.Save(*receiver.Client)
		})
		return
	}

	// 调度失败
	receiver.Task.ScheduleFail(fmt.Sprintf("请求客户端%s（%s）：%s:%d失败:%s", receiver.Client.Name, receiver.Client.Id, receiver.Client.Ip, receiver.Client.Port, err.Error()))
	_ = container.Resolve[redis.IClient]("default").Transaction(func() {
		taskGroupRepository.Save(*receiver.DomainObject)
		clientRepository.Save(*receiver.Client)
	})
}

// 任务完成
func (receiver *TaskGroupMonitor) taskFinish() {
	// 调度失败后，需要立即重新调度
	if receiver.Task.ScheduleStatus != scheduleStatus.Fail && !receiver.Task.IsFinish() {
		return
	}

	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	// 先保存任务内容
	taskGroupRepository.SaveTask(receiver.Task)

	// 计算下一个周期
	if receiver.CalculateNextAtByCron() {
		// 任务初始化
		receiver.CreateTask()
	}
	taskGroupRepository.SaveAndTask(*receiver.DomainObject)
}

// 主动通知客户端，停止任务
func (receiver *TaskGroupMonitor) TaskKill() {
	// FOPS发起Kill请求
	receiver.Client.Kill(receiver.Task.Id)
}

// 通知
func (receiver *TaskGroupMonitor) Notify() {
	flog.Infof("任务组：%s 收到手动更新请求", receiver.Name)
	// 当客户端IsMaster=true，代表当前任务组正在执行，所以才要发消息，否则会导致chan队列撑满
	if receiver.Client != nil && receiver.Client.IsMaster {
		receiver.updated <- struct{}{}
	}
}

// TaskGroupEnableCount 返回开启状态的任务组
func TaskGroupEnableCount() int {
	return taskGroupList.Values().Where(func(item *TaskGroupMonitor) bool {
		return !item.Client.IsClose()
	}).Count()
}

// TaskGroupCount 返回当前正在监控的任务组数量
func TaskGroupCount() int {
	return taskGroupList.Count()
}
