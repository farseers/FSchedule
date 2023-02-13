package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/taskGroup"
	"encoding/json"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
)

// TaskGroupUpdateSubscribe 任务组有更新（Redis订阅）
func TaskGroupUpdateSubscribe(message any, _ core.EventArgs) {
	var taskGroupDO taskGroup.DomainObject
	err := json.Unmarshal([]byte(message.(string)), &taskGroupDO)
	if err != nil {
		return
	}

	flog.Debugf("任务组更新通知：%s Ver%d", taskGroupDO.Name, taskGroupDO.Ver)
	domain.MonitorTaskGroupPush(&taskGroupDO)
}
