package model

import (
	"github.com/farseer-go/collections"
	"time"
)

type TaskGroupPO struct {
	Name        string                                 `gorm:"primaryKey"` // 实现Job的特性名称（客户端识别哪个实现类）
	Ver         int                                    // 版本
	Caption     string                                 // 任务组标题
	StartAt     time.Time                              // 开始时间
	NextAt      time.Time                              // 下次执行时间
	Cron        string                                 // 时间定时器表达式
	ActivateAt  time.Time                              // 活动时间
	LastRunAt   time.Time                              // 最后一次完成时间
	RunSpeedAvg int64                                  // 运行平均耗时
	RunCount    int                                    // 运行次数
	IsEnable    bool                                   // 是否开启
	Data        collections.Dictionary[string, string] `gorm:"serializer:json"` // 传给客户端的参数，按逗号分隔
	Task        TaskPO                                 `gorm:"serializer:json"` // 任务
}
