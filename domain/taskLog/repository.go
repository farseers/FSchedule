package taskLog

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
)

type Repository interface {
	// Add 添加日志
	Add(taskLogDO DomainObject)
	GetList(taskGroupId int64, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.PageList[DomainObject] // 日志列表
}
