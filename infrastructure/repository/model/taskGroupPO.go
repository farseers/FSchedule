package model

import (
	"FSchedule/domain/enum/executeStatus"
	"github.com/farseer-go/collections"
	"time"
)

type TaskGroupPO struct {
	Name              string                                 `gorm:"primaryKey;size:64;not null;comment:任务组名称"`
	Ver               int                                    `gorm:"type:int;not null;comment:版本"`
	Caption           string                                 `gorm:"size:32;not null;comment:任务组标题"`
	StartAt           time.Time                              `gorm:"type:timestamp;size:6;not null;comment:开始时间"`
	NextAt            time.Time                              `gorm:"type:timestamp;size:6;not null;comment:下次执行时间"`
	Cron              string                                 `gorm:"size:32;not null;comment:时间定时器表达式"`
	ActivateAt        time.Time                              `gorm:"type:timestamp;size:6;not null;comment:活动时间"`
	LastRunAt         time.Time                              `gorm:"type:timestamp;size:6;not null;comment:最后一次完成时间"`
	LastExecuteStatus executeStatus.Enum                     `gorm:"type:tinyint;not null;comment:上次执行结果;"`
	RunSpeedAvg       int64                                  `gorm:"type:bigint;not null;comment:运行平均耗时"`
	RunCount          int                                    `gorm:"type:int;not null;comment:运行次数"`
	IsEnable          bool                                   `gorm:"size:1;not null;comment:是否开启"`
	Data              collections.Dictionary[string, string] `gorm:"type:string;size:2048;serializer:json;not null;comment:传给客户端的参数"`
	Task              TaskPO                                 `gorm:"type:string;size:4096;serializer:json;not null;comment:任务"`
	RetryDelaySecond  int                                    `gorm:"type:int;not null;comment:失败后多少秒重试（0不重试）"`
}
