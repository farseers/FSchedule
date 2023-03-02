package schedule

import (
	"github.com/farseer-go/fs/core"
)

type Repository interface {
	// ScheduleLock 创建调度锁
	ScheduleLock(name string, taskId int64) core.ILock
	// Election 选举锁
	Election(fn func())
	// Schedule 调度
	Schedule(name string, fn func())
	// GetLeaderId 获取master集群ID
	GetLeaderId() int64
}
