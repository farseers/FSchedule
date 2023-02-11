package domain

import (
	"FSchedule/domain/client"
	"context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"time"
)

// 客户端列表
var clientList = collections.NewDictionary[int64, *ClientMonitor]()

// ClientMonitor 等待任务执行
type ClientMonitor struct {
	client             *client.DomainObject
	ctx                context.Context
	cancelFunc         context.CancelFunc
	clientRepository   client.Repository
	clientJoinEvent    core.IEvent
	clientOfflineEvent core.IEvent
}

// MonitorClientPush 将最新的客户端信息，推送到监控线程
func MonitorClientPush(clientDO *client.DomainObject) {
	// 新客户端
	if !clientList.ContainsKey(clientDO.Id) {
		ctx, cancelFunc := context.WithCancel(fs.Context)
		clientMonitor := &ClientMonitor{
			client:             clientDO,
			ctx:                ctx,
			cancelFunc:         cancelFunc,
			clientRepository:   container.Resolve[client.Repository](),
			clientJoinEvent:    container.Resolve[core.IEvent]("ClientJoin"),
			clientOfflineEvent: container.Resolve[core.IEvent]("ClientOffline"),
		}
		clientList.Add(clientDO.Id, clientMonitor)

		// 异步检查客户端在线状态
		go clientMonitor.checkOnline()

		// 通知任务组，有新的客户端加入
		_ = clientMonitor.clientJoinEvent.Publish(clientMonitor.client)
	}

	existsClientDO := clientList.GetValue(clientDO.Id)

	// 修改地址对应的值
	*existsClientDO.client = *clientDO

	// 客户端离线
	if existsClientDO.client.IsOffline() {
		_ = existsClientDO.clientOfflineEvent.Publish(existsClientDO.client)
	}
}

// checkOnline 异步检查客户端在线状态
func (receiver *ClientMonitor) checkOnline() {
	for {
		// 离线了，则退出
		if receiver.client.IsOffline() {
			receiver.clientRepository.RemoveClient(receiver.client.Id)
			flog.Infof("客户端（%d）：%s:%d 下线", receiver.client.Id, receiver.client.Ip, receiver.client.Port)
			break
		}

		select {
		case <-time.After(10 * time.Second):
			receiver.client.CheckOnline()
			receiver.clientRepository.Save(receiver.client)

		case <-receiver.ctx.Done():
			break
		}
	}

	// 客户端下线，移除客户端
	clientList.Remove(receiver.client.Id)
}
