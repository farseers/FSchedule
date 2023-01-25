package client

import (
	"FSchedule/domain/enum"
	"time"
)

type DomainObject struct {
	Id         int64                // 客户端ID
	Name       string               // 客户端名称
	Ip         string               // 客户端IP
	Port       int                  // 客户端端口
	ActivateAt time.Time            // 活动时间
	Status     enum.EumClientStatus // 客户端状态
}

// IsNil 判断注册的客户端是否有效
func (receiver DomainObject) IsNil() bool {
	return receiver.Id == 0 || receiver.Name == "" || receiver.Ip == "" || receiver.Port == 0
}
