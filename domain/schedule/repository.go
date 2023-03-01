package schedule

import (
	"github.com/farseer-go/fs/core"
)

type Repository interface {
	// ScheduleLock 创建调度锁
	ScheduleLock(name string) core.ILock
	// Election 选举锁
	Election(fn func())
	// GetLeaderId 获取master集群ID
	GetLeaderId() int64
}
