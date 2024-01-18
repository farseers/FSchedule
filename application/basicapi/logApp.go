// @area /basicapi/
package basicapi

import (
	"FSchedule/domain/taskLog"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
)

// 日志列表
// @get log/list
func LogList(taskGroupName string, logLevel eumLogLevel.Enum, taskId int64, pageSize int, pageIndex int, taskLogRepository taskLog.Repository) collections.PageList[taskLog.DomainObject] {
	return taskLogRepository.GetList(taskGroupName, logLevel, taskId, pageSize, pageIndex)
}

// 日志列表
// @get log/listByClientName
func LogListByClientName(clientName, taskGroupName string, logLevel eumLogLevel.Enum, taskId int64, pageSize int, pageIndex int, taskLogRepository taskLog.Repository) collections.PageList[taskLog.DomainObject] {
	return taskLogRepository.GetListByClientName(clientName, taskGroupName, logLevel, taskId, pageSize, pageIndex)
}
