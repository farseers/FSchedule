package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"encoding/json"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/trace"
)

// ClientUpdateSubscribe 客户端有更新（Redis订阅）
func ClientUpdateSubscribe(message any, _ core.EventArgs) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	var clientDO client.DomainObject
	err := json.Unmarshal([]byte(message.(string)), &clientDO)
	if err != nil || clientDO.IsNil() {
		return
	}
	domain.MonitorClientPush(&clientDO)
}
