package taskLog

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
)

type Repository interface {
	// Add 添加日志
	Add(taskLogDO DomainObject)
	// 获取日志列表
	GetList(taskGroupName string, logLevel eumLogLevel.Enum, taskId int64, pageSize int, pageIndex int) collections.PageList[DomainObject] // 日志列表
}
