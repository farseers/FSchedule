package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/core"
)

// ClientOfflineEvent 客户端下线了
func ClientOfflineEvent(message any, _ core.EventArgs) {
	clientDO := message.(*client.DomainObject)
	domain.ClientOffline(clientDO)
}
