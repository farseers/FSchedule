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

func (repository *TaskLogRepository) GetList(taskGroupId int64, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.PageList[taskLog.DomainObject] {
	ts := context.MysqlContextIns.TaskLog.Desc("create_at")
	if taskGroupId > 0 {
		ts = ts.Where("task_group_id = ?", taskGroupId)
	}
	if logLevel > -1 {
		ts = ts.Where("log_level = ?", logLevel)
	}
	pageList := ts.ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[taskLog.DomainObject](pageList)
}

func (repository *TaskLogRepository) AddBatch(lstPO collections.List[model.TaskLogPO]) {
	_, err := context.MysqlContextIns.TaskLog.InsertList(lstPO, 50)
	if err != nil {
		exception.ThrowRefuseException("批量添加报错：" + err.Error())
	}
}
