package taskGroup

import (
	"FSchedule/domain/enum"
	"time"
)

// ClientVO 客户端
type ClientVO struct {
	Id         int64             // 客户端ID
	Name       string            // 客户端名称
	Ip         string            // 客户端IP
	Port       int               // 客户端端口
	ScheduleAt time.Time         // 任务调度时间
	Status     enum.ClientStatus // 客户端状态
}
