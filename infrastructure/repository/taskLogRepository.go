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

func (repository *TaskLogRepository) GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.PageList[taskLog.DomainObject] {
	pageList := context.MysqlContextIns.TaskLog.Where("name", jobName).Where("log_level", logLevel).ToPageList(pageSize, pageIndex)
	var pageListDO collections.PageList[taskLog.DomainObject]
	pageList.MapToPageList(&pageListDO)
	return pageListDO
}

func (repository *TaskLogRepository) AddBatch(lstPO collections.List[model.TaskLogPO]) {
	err := context.MysqlContextIns.TaskLog.InsertList(lstPO, 50)
	if err != nil {
		exception.ThrowRefuseException("批量添加报错")
	}
}
