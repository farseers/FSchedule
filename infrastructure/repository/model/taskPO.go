package model

import (
	"FSchedule/domain/enum"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/farseer-go/collections"
	"time"
)

type TaskPO struct {
	Id          int64                                  `gorm:"primaryKey;autoIncrement;comment:主键;index:idx_name_status_create,priority:4"`
	Name        string                                 `gorm:"size:64;not null;comment:任务组名称;index:idx_name_create,priority:1;index:idx_name_status_create,priority:1;"`
	Ver         int                                    `gorm:"type:int;not null;comment:版本"`
	Caption     string                                 `gorm:"size:32;not null;comment:任务组标题"`
	StartAt     time.Time                              `gorm:"type:timestamp;size:6;not null;comment:开始时间"`
	RunAt       time.Time                              `gorm:"type:timestamp;size:6;not null;comment:实际执行时间"`
	RunSpeed    int64                                  `gorm:"type:bigint;not null;comment:运行耗时"`
	ClientId    int64                                  `gorm:"type:bigint;not null;comment:客户端Id"`
	ClientIp    string                                 `gorm:"size:32;not null;comment:客户端IP"`
	ClientName  string                                 `gorm:"size:64;not null;comment:客户端名称"`
	Progress    int                                    `gorm:"type:int;not null;comment:进度0-100"`
	Status      enum.TaskStatus                        `gorm:"type:tinyint;not null;comment:状态;index:idx_status_create,priority:1;index:idx_name_status_create,priority:2"`
	SchedulerAt time.Time                              `gorm:"type:timestamp;size:6;not null;comment:调度时间"`
	Data        collections.Dictionary[string, string] `gorm:"type:string;size:2048;serializer:json;not null;comment:本次执行任务时的Data数据"`
	CreateAt    time.Time                              `gorm:"type:timestamp;size:6;not null;comment:任务创建时间;index:idx_status_create,priority:2;index:idx_name_create,priority:2;index:idx_name_status_create,priority:3"`
}

// Value return json value, implement driver.Valuer interface
func (receiver *TaskPO) Value() (driver.Value, error) {
	ba, err := json.Marshal(receiver)
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (receiver *TaskPO) Scan(val any) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	}
	return json.Unmarshal(ba, &receiver)
}

// GormDataType gorm common data type
func (receiver *TaskPO) GormDataType() string {
	return "json"
}
