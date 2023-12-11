// @area /basicapi/
package basicapi

import (
	"FSchedule/domain/taskLog"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
)

// 日志列表
// @get log/list
func LogList(taskGroupId int64, logLevel eumLogLevel.Enum, pageSize int, pageIndex int, taskLogRepository taskLog.Repository) collections.PageList[taskLog.DomainObject] {
	return taskLogRepository.GetList(taskGroupId, logLevel, pageSize, pageIndex)
}
