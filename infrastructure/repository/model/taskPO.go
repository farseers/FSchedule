package model

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"time"
)

type TaskPO struct {
	Id             int64                                  `gorm:"primaryKey;autoIncrement;comment:主键;"`
	TraceId        string                                 `gorm:"not null;default:'';comment:上下文ID"`
	Name           string                                 `gorm:"size:64;not null;comment:任务组名称;"`
	Ver            int                                    `gorm:"type:int;not null;comment:版本"`
	Caption        string                                 `gorm:"size:32;not null;comment:任务组标题"`
	StartAt        time.Time                              `gorm:"type:timestamp;size:6;not null;comment:开始时间"`
	RunAt          time.Time                              `gorm:"type:timestamp;size:6;not null;comment:实际执行时间"`
	RunSpeed       int64                                  `gorm:"type:bigint;not null;comment:运行耗时"`
	ClientId       int64                                  `gorm:"type:bigint;not null;comment:客户端Id"`
	ClientIp       string                                 `gorm:"size:32;not null;comment:客户端IP"`
	ClientName     string                                 `gorm:"size:64;not null;comment:客户端名称"`
	Progress       int                                    `gorm:"type:int;not null;comment:进度0-100"`
	ExecuteStatus  executeStatus.Enum                     `gorm:"type:tinyint;not null;comment:执行状态;"`
	ScheduleStatus scheduleStatus.Enum                    `gorm:"type:tinyint;not null;comment:调度状态;"`
	SchedulerAt    time.Time                              `gorm:"type:timestamp;size:6;not null;comment:调度时间"`
	Data           collections.Dictionary[string, string] `gorm:"type:string;size:2048;json;not null;comment:本次执行任务时的Data数据"`
	CreateAt       time.Time                              `gorm:"type:timestamp;size:6;not null;comment:任务创建时间;"`
	Remark         string                                 `gorm:"size:1024;not null;comment:备注"`
}

// 创建索引
func (*TaskPO) CreateIndex() map[string]data.IdxField {
	return map[string]data.IdxField{
		/*
			Where("name = ?", taskGroupName)
			name = ? and (execute_status = ? or execute_status = ?)
		*/
		"idx_name_create": {false, "name,create_at desc"},
		/*
			execute_status = ? and create_at >= DATE_SUB(CURDATE(), INTERVAL 3 DAY) group by name
			name = ? and (execute_status = ? or execute_status = ?) and create_at < ? and ClientId < ?
			(execute_status = ? or execute_status = ?) and (create_at >= ?)
			execute_status = ? and create_at >= ?
		*/
		"idx_name_status_create": {false, "create_at desc,execute_status,name"},
	}
}
