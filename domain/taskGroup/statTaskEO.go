package taskGroup

import (
	"FSchedule/domain/enum"
)

type StatTaskEO struct {
	ClientName string          // 应用名称
	Status     enum.TaskStatus // 状态
	Count      int             // 日志数量
}
