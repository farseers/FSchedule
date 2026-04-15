package schedule

import (
	"context"

	"github.com/farseer-go/fs/core"
)

type Repository interface {
	// RegistryLock 创建注册锁
	RegistryLock(clientId int64) core.ILock
	// Election 选举锁
	Election(ctx context.Context, fn func())
	// 监控任务组超时锁
	Monitor(ctx context.Context, fn func())
	// Schedule 调度
	Schedule(ctx context.Context, taskGroupName string, fn func())
	// GetLeaderId 获取master集群ID
	GetLeaderId() int64
}
