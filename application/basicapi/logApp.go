// @area /basicapi/
package basicapi

import (
	"FSchedule/domain/taskLog"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
)

// 日志列表
// @get log/list
func LogList(taskGroupName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int, taskLogRepository taskLog.Repository) collections.PageList[taskLog.DomainObject] {
	return taskLogRepository.GetList(taskGroupName, logLevel, pageSize, pageIndex)
}
