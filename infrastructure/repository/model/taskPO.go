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
	Id          int64                                  `gorm:"primaryKey"` // 主键
	Name        string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver         int                                    // 版本
	Caption     string                                 // 任务组标题
	StartAt     time.Time                              // 开始时间
	RunAt       time.Time                              // 实际执行时间
	RunSpeed    int64                                  // 运行耗时
	ClientId    int64                                  // 客户端Id
	ClientIp    string                                 // 客户端IP
	ClientName  string                                 // 客户端名称
	Progress    int                                    // 进度0-100
	Status      enum.TaskStatus                        // 状态
	SchedulerAt time.Time                              // 调度时间
	Data        collections.Dictionary[string, string] `gorm:"serializer:json"` // 本次执行任务时的Data数据
	CreateAt    time.Time                              // 任务创建时间
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
