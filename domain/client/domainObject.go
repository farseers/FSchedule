package client

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"time"
)

type DomainObject struct {
	Id          int64                   // 客户端ID
	Name        string                  // 客户端名称
	Ip          string                  // 客户端IP
	Port        int                     // 客户端端口
	ActivateAt  dateTime.DateTime       // 活动时间
	ScheduleAt  dateTime.DateTime       // 任务调度时间
	Status      enum.ClientStatus       // 客户端状态
	QueueCount  int                     // 排队中的任务数量
	WorkCount   int                     // 正在处理的任务数量
	CpuUsage    float64                 // CPU百分比
	MemoryUsage float64                 // 内存百分比
	ErrorCount  int                     // 错误次数
	Jobs        collections.List[JobVO] // 客户端支持的任务
	NeedNotice  bool                    //	是否需要通知任务组
}

// IsNil 判断注册的客户端是否有效
func (receiver *DomainObject) IsNil() bool {
	return receiver.Id == 0 || receiver.Name == "" || receiver.Ip == "" || receiver.Port == 0
}

// IsOnline 是否刚注册进来
func (receiver *DomainObject) IsOnline() bool {
	return receiver.Status == enum.Online
}

// IsOffline 判断客户端是否下线
func (receiver *DomainObject) IsOffline() bool {
	return receiver.Status == enum.Offline
}

// IsNotSchedule 状态不是调度状态
func (receiver *DomainObject) IsNotSchedule() bool {
	return receiver.Status != enum.Scheduler
}

// Registry 注册客户端
func (receiver *DomainObject) Registry() {
	receiver.ActivateAt = dateTime.Now()
	receiver.Status = enum.Online
	receiver.NeedNotice = true
}

// Logout 客户端下线
func (receiver *DomainObject) Logout() {
	receiver.Status = enum.Offline
	receiver.NeedNotice = true
}

// CheckOnline 检查客户端是否存活
func (receiver *DomainObject) CheckOnline() error {
	status, err := container.Resolve[IClientCheck]().Check(receiver)
	if err != nil {
		flog.Warningf("检查客户端%s（%d）：%s:%d 是否存活失败：%s", receiver.Name, receiver.Id, receiver.Ip, receiver.Port, err.Error())
	}
	receiver.updateStatus(status, err)
	return err
}

// Schedule 调度
func (receiver *DomainObject) Schedule(task TaskEO) bool {
	status, err := container.Resolve[IClientCheck]().Invoke(receiver, task)
	flog.Warningf("任务组：%s %d 向客户端%s（%d）：%s:%d 调度失败：%s", task.Name, task.Id, receiver.Name, receiver.Id, receiver.Ip, receiver.Port, err.Error())
	receiver.updateStatus(status, err)

	milliseconds := time.Since(task.StartAt).Milliseconds()
	if milliseconds < 0 {
		milliseconds = 0
	}
	if receiver.Status == enum.Scheduler {
		receiver.ScheduleAt = dateTime.Now()

		//flog.Infof("任务组：%s %d 调度成功 延迟：%s ms", task.Name, task.Id, flog.Red(milliseconds))
		return true
	}
	return false
}

// 更新状态
func (receiver *DomainObject) updateStatus(status ResourceVO, err error) {
	oldStatus := receiver.Status
	if err != nil {
		// 先设置为无法调度
		receiver.UnSchedule()
	} else {
		receiver.ActivateAt = dateTime.Now()
		receiver.ErrorCount = 0
		receiver.CpuUsage = status.CpuUsage
		receiver.MemoryUsage = status.MemoryUsage
		receiver.QueueCount = status.QueueCount
		receiver.WorkCount = status.WorkCount

		if status.AllowSchedule {
			receiver.Status = enum.Scheduler
		} else {
			receiver.Status = enum.StopSchedule
		}
	}

	receiver.NeedNotice = oldStatus != receiver.Status
}

// UnSchedule 客户端无法调度
func (receiver *DomainObject) UnSchedule() {
	if !receiver.IsOffline() {
		receiver.ErrorCount++
		receiver.Status = enum.UnSchedule

		// 大于3次、活动时间超过30秒，则判定为离线
		now := dateTime.Now()
		if receiver.ErrorCount >= 3 && now.Sub(receiver.ActivateAt).Seconds() >= 30 {
			receiver.Logout()
		}
	}
}
