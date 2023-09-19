package client

import (
	"github.com/farseer-go/collections"
	"time"
)

// TaskEO 任务记录
type TaskEO struct {
	Id          int64                                  // 主键
	Caption     string                                 // 任务组标题
	TaskGroupId int64                                  // 任务组ID
	Name        string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	StartAt     time.Time                              // 开始时间
	Data        collections.Dictionary[string, string] // 本次执行任务时的Data数据
}
