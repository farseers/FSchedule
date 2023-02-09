package client

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"time"
)

type DomainObject struct {
	Id          int64                   // 客户端ID
	Name        string                  // 客户端名称
	Ip          string                  // 客户端IP
	Port        int                     // 客户端端口
	ActivateAt  time.Time               // 活动时间
	ScheduleAt  time.Time               // 任务调度时间
	Status      enum.ClientStatus       // 客户端状态
	QueueCount  int                     // 排队中的任务数量
	WorkCount   int                     // 正在处理的任务数量
	CpuUsage    float32                 // CPU百分比
	MemoryUsage float32                 // 内存百分比
	ErrorCount  int                     // 错误次数
	Jobs        collections.List[JobVO] // 客户端支持的任务
}

// IsNil 判断注册的客户端是否有效
func (receiver *DomainObject) IsNil() bool {
	return receiver.Id == 0 || receiver.Name == "" || receiver.Ip == "" || receiver.Port == 0
}

// IsOffline 判断客户端是否下线
func (receiver *DomainObject) IsOffline() bool {
	return receiver.Status == enum.Offline
}

// Registry 注册客户端
func (receiver *DomainObject) Registry() {
	receiver.ActivateAt = time.Now()
	receiver.Status = enum.Online
}

// CheckOnline 检查客户端是否存活
func (receiver *DomainObject) CheckOnline() {
	if !receiver.IsOffline() {
		status, err := container.Resolve[IClientCheck]().Check(receiver)
		receiver.updateStatus(status, err)
	}
}

// Logout 客户端下线
func (receiver *DomainObject) Logout() {
	receiver.Status = enum.Offline
}

// Schedule 调度
func (receiver *DomainObject) Schedule(task *TaskEO) bool {
	status, err := container.Resolve[IClientCheck]().Invoke(receiver, task)
	receiver.updateStatus(status, err)

	if receiver.Status == enum.Scheduler {
		receiver.ScheduleAt = time.Now()
	}

	return receiver.Status == enum.Scheduler
}

// 更新状态
func (receiver *DomainObject) updateStatus(status ResourceVO, err error) {
	if err != nil {
		// 先设置为无法调度
		receiver.UnSchedule()
	} else {
		receiver.ActivateAt = time.Now()
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
}

// UnSchedule 客户端无法调度
func (receiver *DomainObject) UnSchedule() {
	receiver.ErrorCount++
	receiver.Status = enum.UnSchedule

	// 大于3次、活动时间超过30秒，则判定为离线
	if receiver.ErrorCount > 3 && time.Now().Sub(receiver.ActivateAt).Seconds() >= 30 {
		receiver.Logout()
	}
}
