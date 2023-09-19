package taskGroupApp

import (
	"FSchedule/domain/taskGroup"
	"FSchedule/domain/taskLog"
	"github.com/farseer-go/fs/core/eumLogLevel"
)

type logReportDTO struct {
	Logs []LogContent
}

type LogContent struct {
	TaskId      int64  // 主键
	TaskGroupId int64  // 任务组ID
	Name        string // 实现Job的特性名称（客户端识别哪个实现类）
	LogLevel    eumLogLevel.Enum
	CreateAt    int64
	Content     string
}

// LogReport 日志上报
func LogReport(dto logReportDTO, taskGroupRepository taskGroup.Repository, taskLogRepository taskLog.Repository) {
	for _, log := range dto.Logs {
		taskDO := taskGroupRepository.GetTask(log.TaskId, log.TaskId)
		taskLogDO := taskLog.NewDO(log.Name, taskDO.Caption, taskDO.Ver, taskDO.Id, taskDO.Data, log.LogLevel, log.Content, log.CreateAt)
		taskLogRepository.Add(taskLogDO)
	}
}
