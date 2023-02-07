package client

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"time"
)

type DomainObject struct {
	Id          int64             // 客户端ID
	Name        string            // 客户端名称
	Ip          string            // 客户端IP
	Port        int               // 客户端端口
	ActivateAt  time.Time         // 活动时间
	ScheduleAt  time.Time         // 任务调度时间
	Status      enum.ClientStatus // 客户端状态
	QueueCount  int               // 排队中的任务数量
	WorkCount   int               // 正在处理的任务数量
	CpuUsage    float32           // CPU百分比
	MemoryUsage float32           // 内存百分比
	Jobs        []JobVO           // 客户端支持的任务
}

// IsNil 判断注册的客户端是否有效
func (receiver *DomainObject) IsNil() bool {
	return receiver.Id == 0 || receiver.Name == "" || receiver.Ip == "" || receiver.Port == 0
}

// Registry 注册客户端
func (receiver *DomainObject) Registry() {
	receiver.ActivateAt = time.Now()
	receiver.Status = enum.Online
}

// CheckOnline 检查客户端是否存活
func (receiver *DomainObject) CheckOnline() {
	// 只检查非离线状态
	if receiver.Status != enum.Offline {
		status, err := container.Resolve[IClientCheck]().Check(receiver)
		receiver.updateStatus(status, err)
	}

	if receiver.Status == enum.Offline {
		receiver.Logout()
	}
}

// Logout 客户端下线
func (receiver *DomainObject) Logout() {
	container.Resolve[core.IEvent]("ClientOffline").Publish(receiver)
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
		receiver.Status = enum.UnSchedule
		// 如果活动时间超过30秒，则判定为离线状态
		if time.Now().Sub(receiver.ActivateAt).Seconds() >= 30 {
			receiver.Status = enum.Offline
		}
	} else {
		receiver.ActivateAt = time.Now()
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
