package model

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"time"
)

type TaskLogPO struct {
	Id       int64                                  `gorm:"primaryKey"` // 主键
	Name     string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Caption  string                                 // 任务组标题
	Ver      int                                    // 版本
	TaskId   int64                                  // 任务ID
	Data     collections.Dictionary[string, string] `gorm:"serializer:json"` // 本次执行任务时的Data数据
	LogLevel eumLogLevel.Enum                       // 日志级别
	Content  string                                 // 日志内容
	CreateAt time.Time                              // 日志时间
}
