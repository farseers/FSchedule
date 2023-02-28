package repository

import (
	"FSchedule/domain/taskLog"
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/queue"
)

type TaskLogRepository struct {
	TaskLog data.TableSet[model.TaskLogPO] `data:"name=task_log"`
}

func (repository TaskLogRepository) Add(taskLogDO taskLog.DomainObject) {
	po := mapper.Single[model.TaskLogPO](taskLogDO)
	queue.Push("TaskLogQueue", po)
}

func (repository TaskLogRepository) GetList(jobName string, logLevel eumLogLevel.Enum, pageSize int, pageIndex int) collections.PageList[taskLog.DomainObject] {
	pageList := repository.TaskLog.Where("name", jobName).Where("log_level", logLevel).ToPageList(pageSize, pageIndex)
	var pageListDO collections.PageList[taskLog.DomainObject]
	pageList.MapToPageList(&pageListDO)
	return pageListDO
}

func (repository TaskLogRepository) AddBatch(lstPO collections.List[model.TaskLogPO]) {
	err := repository.TaskLog.InsertList(lstPO, 50)
	if err != nil {
		exception.ThrowRefuseException("批量添加报错")
	}
}
