package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/timingWheel"
	"time"
)

// 客户端列表
var clientList = collections.NewDictionary[int64, *ClientMonitor]()

// ClientMonitor 等待任务执行
type ClientMonitor struct {
	client           *client.DomainObject
	ctx              context.Context
	cancelFunc       context.CancelFunc
	ClientRepository client.Repository
}

// MonitorClientPush 将最新的客户端信息，推送到监控线程
func MonitorClientPush(clientDO *client.DomainObject) {
	// 新客户端
	if !clientDO.IsOffline() && !clientList.ContainsKey(clientDO.Id) {
		ctx, cancelFunc := context.WithCancel(fs.Context)
		clientMonitor := container.ResolveIns(&ClientMonitor{
			client:     clientDO,
			ctx:        ctx,
			cancelFunc: cancelFunc,
		})
		clientList.Add(clientDO.Id, clientMonitor)

		flog.Infof("客户端（%d）开始监听：%s:%d", clientDO.Id, clientDO.Ip, clientDO.Port)

		// 异步检查客户端在线状态
		go clientMonitor.checkOnline()

		ClientUpdate(clientDO)
	}

	existsClientDO := clientList.GetValue(clientDO.Id)

	if existsClientDO != nil {
		// 修改地址对应的值
		*existsClientDO.client = *clientDO

		// 只有状态不一样时，才要更新
		if existsClientDO.client.IsNotSchedule() {
			ClientUpdate(existsClientDO.client)
		}

		// 客户端离线
		if clientDO.IsOffline() {
			flog.Infof("客户端（%d）：%s:%d 下线", clientDO.Id, clientDO.Ip, clientDO.Port)
			existsClientDO.cancelFunc()
			clientList.Remove(clientDO.Id)
		}
	}
}

// checkOnline 异步检查客户端在线状态
func (receiver *ClientMonitor) checkOnline() {
	for {
		if receiver.client.IsOffline() {
			return
		}
		checkTime := 60 * time.Second
		if receiver.client.IsNotSchedule() || time.Since(receiver.client.ActivateAt).Seconds() >= 60 {
			checkTime = 10 * time.Second
		}

		select {
		case <-timingWheel.Add(checkTime).C:
			if !receiver.client.IsOffline() {
				receiver.client.CheckOnline()
				receiver.ClientRepository.Save(receiver.client)
			}
		case <-receiver.ctx.Done():
			return
		}
	}
}

// ClientCount 返回当前正在监控的客户端数量
func ClientCount() int {
	return clientList.Count()
}

// ClientNormalCount 返回正常状态的客户端数量
func ClientNormalCount() int {
	return clientList.Values().Where(func(item *ClientMonitor) bool {
		return item.client.Status == enum.Scheduler
	}).Count()
}
