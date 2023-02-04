package domainEvent

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
)

// RemoveClientEvent 客户端离线
func RemoveClientEvent(message any, _ eventBus.EventArgs) {
	do := message.(*client.DomainObject)
	// 仓储移除客户端
	repository := container.Resolve[client.Repository]()
	repository.RemoveClient(do.Id)

	flog.Infof("客户端：%s:%d %d 下线", do.Ip, do.Port, do.Id)
}
