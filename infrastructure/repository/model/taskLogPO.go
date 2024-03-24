package model

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"time"
)

type TaskLogPO struct {
	Id       int64                                  `gorm:"primaryKey;autoIncrement;comment:主键"`
	Name     string                                 `gorm:"size:64;not null;comment:任务组名称;"`
	Ver      int                                    `gorm:"type:int;not null;comment:版本"`
	Caption  string                                 `gorm:"size:32;not null;comment:任务组标题"`
	TaskId   int64                                  `gorm:"type:bigint;not null;comment:任务ID"`
	Data     collections.Dictionary[string, string] `gorm:"type:string;size:2048;serializer:json;not null;comment:本次执行任务时的Data数据"`
	LogLevel eumLogLevel.Enum                       `gorm:"type:tinyint;not null;comment:日志级别;"`
	Content  string                                 `gorm:"type:text;size:0;not null;comment:日志内容"`
	CreateAt time.Time                              `gorm:"type:timestamp;size:6;not null;comment:日志时间;"`
}

// 创建索引
func (*TaskLogPO) CreateIndex() map[string]data.IdxField {
	return map[string]data.IdxField{
		"idx_name_logLevel": {false, "name,log_level,create_at desc"},
		"idx_name":          {false, "name,create_at desc"},
		"idx_create_at":     {false, "create_at desc"},
	}
}
