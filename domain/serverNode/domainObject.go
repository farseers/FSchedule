package serverNode

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/parse"
	"strings"
)

var IsLeaderNode bool

type DomainObject struct {
	Id         int64             // 客户端ID
	Name       string            // 客户端名称
	Ip         string            // 客户端IP
	Port       int               // 客户端端口
	IsLeader   bool              // 是否为Master
	ActivateAt dateTime.DateTime // 活动时间
}

func New() *DomainObject {
	addr := configure.GetString("WebApi.Url")
	if addr == "" {
		addr = ":8888"
	}
	addr, _ = strings.CutPrefix(addr, ":")
	return &DomainObject{
		Id:         core.AppId,
		Name:       core.HostName,
		Ip:         core.AppIp,
		Port:       parse.Convert(addr, 0),
		ActivateAt: dateTime.Now(),
	}
}

// SetLeader 设为master
func (receiver *DomainObject) SetLeader(leaderId int64) {
	receiver.IsLeader = leaderId == receiver.Id
}

// Activate 更新活跃时间
func (receiver *DomainObject) Activate() {
	receiver.ActivateAt = dateTime.Now()
}
