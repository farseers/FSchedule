package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"
	"encoding/json"
	"github.com/farseer-go/fs/core"
)

// TaskGroupUpdateSubscribe 任务组有更新（Redis订阅）
func TaskGroupUpdateSubscribe(message any, _ core.EventArgs) {
	var taskGroupDO taskGroup.DomainObject
	err := json.Unmarshal([]byte(message.(string)), &taskGroupDO)
	if err != nil {
		return
	}

	domain.MonitorTaskGroupPush(&taskGroupDO)
}
