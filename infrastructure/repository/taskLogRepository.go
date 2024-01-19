package repository

import (
	"FSchedule/domain/taskLog"
	"FSchedule/infrastructure/repository/context"
	"FSchedule/infrastructure/repository/model"
	"bytes"
	_ "embed"
	"fmt"
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
	ts := context.MysqlContextIns("获取日志列表").TaskLog.Desc("create_at").
		WhereIf(taskGroupName != "", "name = ?", taskGroupName).
		WhereIf(logLevel > -1, "log_level >= ?", logLevel).
		WhereIf(taskId > 0, "task_id = ?", taskId)

	pageList := ts.Desc("create_at").ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[taskLog.DomainObject](pageList)
}

func (repository *TaskLogRepository) AddBatch(lstPO collections.List[model.TaskLogPO]) {
	_, err := context.MysqlContextIns("批量添加日志").TaskLog.InsertList(lstPO, 50)
	if err != nil {
		exception.ThrowRefuseException("批量添加报错：" + err.Error())
	}
}

func (repository *TaskLogRepository) DeleteLog(taskGroupName string) {
	_, _ = context.MysqlContextIns("删除日志").TaskLog.Where("name = ?", taskGroupName).Delete()
}

//go:embed model/sql/taskJoinTaskLog.sql
var taskJoinTaskLogSql string

//go:embed model/sql/taskJoinTaskLogCount.sql
var taskJoinTaskLogCountSql string

func (repository *TaskLogRepository) GetListByClientName(clientName, taskGroupName string, logLevel eumLogLevel.Enum, taskId int64, pageSize int, pageIndex int) collections.PageList[taskLog.DomainObject] {
	var where bytes.Buffer
	if len(clientName) > 0 {
		where.WriteString(fmt.Sprintf(" and task.client_name ='%s'", clientName))
	}
	if len(taskGroupName) > 0 {
		where.WriteString(fmt.Sprintf(" and log.name ='%s'", taskGroupName))
	}
	if logLevel > -1 {
		where.WriteString(fmt.Sprintf(" and log.log_level >= %d", logLevel))
	}
	if taskId > 0 {
		where.WriteString(fmt.Sprintf(" and log.task_id = %d", taskId))
	}
	lst := context.MysqlContextIns("获取日志列表").TaskLog.ExecuteSqlToList(fmt.Sprintf(taskJoinTaskLogSql, where.String(), (pageIndex-1)*pageSize, pageSize))
	count := int64(0)
	_, _ = context.MysqlContextIns("获取日志列表数量").ExecuteSqlToValue(&count, fmt.Sprintf(taskJoinTaskLogCountSql, where.String()))
	return collections.NewPageList[taskLog.DomainObject](mapper.ToList[taskLog.DomainObject](lst), count)
}
