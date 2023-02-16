package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
)

// ClientOfflineEvent 客户端下线了
func ClientOfflineEvent(message any, _ core.EventArgs) {
	clientDO := message.(*client.DomainObject)
	container.Resolve[client.Repository]().RemoveClient(clientDO.Id)
	flog.Infof("客户端（%d）：%s:%d 下线", clientDO.Id, clientDO.Ip, clientDO.Port)
	domain.ClientUpdate(clientDO)
}
