package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum/clientStatus"
	"FSchedule/domain/serverNode"
	"context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/timingWheel"
	"github.com/farseer-go/fs/trace"
	"time"
)

// 客户端列表
type ClientId = int64

var clientList = collections.NewDictionary[ClientId, *ClientMonitor]()

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
		ctx, cancelFunc := context.WithCancel(context.Background())
		clientMonitor := container.ResolveIns(&ClientMonitor{
			client:     clientDO,
			ctx:        ctx,
			cancelFunc: cancelFunc,
		})
		clientList.Add(clientDO.Id, clientMonitor)

		flog.Infof("客户端%s（%d）开始监听：%s:%d", clientDO.Name, clientDO.Id, clientDO.Ip, clientDO.Port)

		// 异步检查客户端在线状态
		if serverNode.IsLeaderNode {
			go clientMonitor.checkOnline()
		}

		//ClientUpdate(clientDO)
	}

	existsClientDO := clientList.GetValue(clientDO.Id)
	if existsClientDO != nil {
		// 修改地址对应的值
		*existsClientDO.client = *clientDO
		ClientUpdate(existsClientDO.client)

		// 客户端离线
		if clientDO.IsOffline() {
			flog.Infof("客户端%s（%d）：%s:%d 下线", clientDO.Name, clientDO.Id, clientDO.Ip, clientDO.Port)
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
		// 不可调度状态，则10秒后检查
		if receiver.client.IsNotSchedule() || dateTime.Since(receiver.client.ActivateAt).Seconds() >= 60 {
			checkTime = 10 * time.Second
		}
		// 新注册，则在3秒后立即检查
		if receiver.client.IsOnline() {
			checkTime = 3 * time.Second
		}

		select {
		case <-timingWheel.Add(checkTime).C:
			// 客户端接受调度状态，且60秒内有活动的，不需要检查
			if receiver.client.Status == clientStatus.Scheduler && dateTime.Now().Sub(receiver.client.ActivateAt).Seconds() < 60 {
				continue
			}

			// 检查非离线状态
			if !receiver.client.IsOffline() {
				// 链路追踪
				traceContext := container.Resolve[trace.IManager]().EntryTask("检查客户端在线状态")

				err := receiver.client.CheckOnline()
				receiver.ClientRepository.Save(receiver.client)
				traceContext.Error(err)
				traceContext.End()
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
		return item.client.Status == clientStatus.Scheduler
	}).Count()
}

func CheckOnline() {
	clients := clientList.Values()
	for i := 0; i < clients.Count(); i++ {
		go clients.Index(i).checkOnline()
	}
}
