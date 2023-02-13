package localQueue

import (
	"FSchedule/domain/taskLog"
	"FSchedule/infrastructure/repository"
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
)

// TaskLogQueueConsumer 将日志指写入
func TaskLogQueueConsumer(subscribeName string, message collections.ListAny, remainingCount int) {
	// 转成BuildLogVO数组
	var lstPO collections.List[model.TaskLogPO]
	message.MapToList(&lstPO)
	container.Resolve[taskLog.Repository]().(repository.TaskLogRepository).AddBatch(lstPO)
}
