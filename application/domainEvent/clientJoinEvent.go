package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/core"
)

// ClientJoinEvent 新的客户端加入
func ClientJoinEvent(message any, _ core.EventArgs) {
	clientDO := message.(*client.DomainObject)
	domain.ClientJoin(clientDO)
}
