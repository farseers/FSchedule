package model

import (
	"FSchedule/domain/enum"
	"time"
)

type TaskPO struct {
	Id          int             `gorm:"primaryKey"` // 主键
	TaskGroupId int             // 任务组ID
	Caption     string          // 任务组标题
	Name        string          // 实现Job的特性名称（客户端识别哪个实现类）
	StartAt     time.Time       // 开始时间
	RunAt       time.Time       // 实际执行时间
	RunSpeed    int64           // 运行耗时
	ClientId    int64           // 客户端Id
	ClientIp    string          // 客户端IP
	ClientName  string          // 客户端名称
	Progress    int             // 进度0-100
	Status      enum.TaskStatus // 状态
	CreateAt    time.Time       // 任务创建时间
	SchedulerAt time.Time       // 调度时间
}
