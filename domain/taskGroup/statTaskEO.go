package taskGroup

import (
	"FSchedule/domain/enum/executeStatus"
)

type StatTaskEO struct {
	ClientName    string             // 应用名称
	ExecuteStatus executeStatus.Enum // 状态
	Count         int                // 日志数量
}
