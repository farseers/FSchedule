package repository

import (
	"FSchedule/domain/taskLog"
	"FSchedule/infrastructure/repository/context"
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/queue"
)

type TaskLogRepository struct {
}

func (repository *TaskLogRepository) Add(taskLogDO taskLog.DomainObject) {
	po := mapper.Single[model.TaskLogPO](taskLogDO)
	queue.Push("TaskLogQueue", po)
}

func (repository *TaskLogRepository) GetList(taskGroupName string, logLevel eumLogLevel.Enum, taskId int64, pageSize int, pageIndex int) collections.PageList[taskLog.DomainObject] {
	ts := context.MysqlContextIns.TaskLog.Desc("create_at").
		WhereIf(taskGroupName != "", "name = ?", taskGroupName).
		WhereIf(logLevel > -1, "log_level >= ?", logLevel).
		WhereIf(taskId > 0, "task_id = ?", taskId)

	pageList := ts.Desc("create_at").ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[taskLog.DomainObject](pageList)
}

func (repository *TaskLogRepository) AddBatch(lstPO collections.List[model.TaskLogPO]) {
	_, err := context.MysqlContextIns.TaskLog.InsertList(lstPO, 50)
	if err != nil {
		exception.ThrowRefuseException("批量添加报错：" + err.Error())
	}
}

func (repository *TaskLogRepository) DeleteLog(taskGroupName string) {
	_, _ = context.MysqlContextIns.TaskLog.Where("name = ?", taskGroupName).Delete()
}
