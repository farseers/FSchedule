package model

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"time"
)

type TaskLogPO struct {
	Id       int64                                  `gorm:"primaryKey;autoIncrement;comment:主键"`
	Name     string                                 `gorm:"size:64;not null;index:idx_name_logLevel,priority:1;comment:任务组名称"`
	Ver      int                                    `gorm:"type:int;not null;comment:版本"`
	Caption  string                                 `gorm:"size:32;not null;comment:任务组标题"`
	TaskId   int64                                  `gorm:"type:bigint;not null;comment:任务ID"`
	Data     collections.Dictionary[string, string] `gorm:"type:string;size:2048;serializer:json;not null;comment:本次执行任务时的Data数据"`
	LogLevel eumLogLevel.Enum                       `gorm:"type:tinyint;not null;index:idx_name_logLevel,priority:2;comment:日志级别"`
	Content  string                                 `gorm:"type:text;size:0;not null;comment:日志内容"`
	CreateAt time.Time                              `gorm:"type:timestamp;size:6;not null;comment:日志时间"`
}
