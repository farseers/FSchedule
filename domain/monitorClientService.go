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
	//flog.Debugf("客户端（%d）更新通知：%s:%d", clientDO.Id, clientDO.Ip, clientDO.Port)
	// 新客户端
	if !clientDO.IsOffline() && !clientList.ContainsKey(clientDO.Id) {
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

	if existsClientDO != nil {
		// 修改地址对应的值
		*existsClientDO.client = *clientDO
		ClientUpdate(existsClientDO.client)
	}

	// 客户端离线
	if clientDO.IsOffline() {
		if existsClientDO != nil {
			existsClientDO.cancelFunc()
			// 客户端下线，移除客户端
			_ = existsClientDO.clientOfflineEvent.Publish(clientDO)
			clientList.Remove(clientDO.Id)
		}
		container.Resolve[client.Repository]().RemoveClient(clientDO.Id)
		flog.Infof("客户端（%d）：%s:%d 下线", clientDO.Id, clientDO.Ip, clientDO.Port)
	}
}

// checkOnline 异步检查客户端在线状态
func (receiver *ClientMonitor) checkOnline() {
	flog.Infof("客户端（%d）开始监听：%s:%d", receiver.client.Id, receiver.client.Ip, receiver.client.Port)
	for {
		if receiver.client.IsOffline() {
			return
		}
		select {
		case <-time.After(10 * time.Second):
			receiver.client.CheckOnline()
			receiver.clientRepository.Save(receiver.client)
		case <-receiver.ctx.Done():
			return
		}
	}
}

// ClientCount 返回当前正在监控的客户端数量
func ClientCount() int {
	return clientList.Count()
}
