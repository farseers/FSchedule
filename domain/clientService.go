package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/schedule"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/timingWheel"
	"github.com/farseer-go/webapi/websocket"
	"time"
)

// 每5秒自动激活客户端的活动时间
func ActivateClient(ctx *websocket.BaseContext, clientId int64, clientRepository client.Repository, scheduleRepository schedule.Repository) {
	// 由于创建锁的时候，需要网络IO开销，所以这里提前100ms进入
	for {
		<-timingWheel.Add(5 * time.Second).C
		// 客户端断开，则退出
		if ctx.IsClose() {
			return
		}
		scheduleRepository.RegistryLock(clientId).GetLockRun(func() {
			clientDO := clientRepository.ToEntity(clientId)
			clientDO.ActivateAt = dateTime.Now()
			clientRepository.Save(clientDO)
		})
	}
}
