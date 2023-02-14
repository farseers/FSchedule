package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"encoding/json"
	"github.com/farseer-go/fs/core"
)

// ClientUpdateSubscribe 客户端有更新（Redis订阅）
func ClientUpdateSubscribe(message any, _ core.EventArgs) {
	var clientDO client.DomainObject
	err := json.Unmarshal([]byte(message.(string)), &clientDO)
	if err != nil {
		return
	}
	domain.MonitorClientPush(&clientDO)
}
