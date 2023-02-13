package taskGroupApp

import (
	"FSchedule/domain/taskGroup"
	"FSchedule/domain/taskLog"
	"github.com/farseer-go/fs/core/eumLogLevel"
)

type logReportDTO struct {
	TaskId int64  // 主键
	Name   string // 实现Job的特性名称（客户端识别哪个实现类）
	Log    []LogContent
}

type LogContent struct {
	LogLevel eumLogLevel.Enum
	CreateAt int64
	Content  string
}

// LogReport 日志上报
func LogReport(dto logReportDTO, taskGroupRepository taskGroup.Repository, taskLogRepository taskLog.Repository) {
	taskDO := taskGroupRepository.GetTask(dto.Name, dto.TaskId)
	for _, log := range dto.Log {
		taskLogDO := taskLog.NewDO(dto.Name, taskDO.Caption, taskDO.Ver, taskDO.Id, taskDO.Data, log.LogLevel, log.Content, log.CreateAt)
		taskLogRepository.Add(taskLogDO)
	}
}
