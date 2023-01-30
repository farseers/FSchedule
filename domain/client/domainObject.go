package client

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"time"
)

type DomainObject struct {
	Id         int64             // 客户端ID
	Name       string            // 客户端名称
	Ip         string            // 客户端IP
	Port       int               // 客户端端口
	ActivateAt time.Time         // 活动时间
	Status     enum.ClientStatus // 客户端状态
	Jobs       []JobVO           // 客户端支持的任务
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
		status := container.Resolve[IClientCheck]().Check(receiver)
		if status {
			receiver.ActivateAt = time.Now()
			receiver.Status = enum.Scheduling
		} else {
			if time.Now().Sub(receiver.ActivateAt).Seconds() >= 30 {
				receiver.Status = enum.Offline
			}
		}
	}

	if receiver.Status == enum.Offline {
		receiver.Logout()
	}
}

// Logout 客户端下线
func (receiver *DomainObject) Logout() {
	container.Resolve[core.IEvent]("ClientOffline").Publish(receiver)
}
