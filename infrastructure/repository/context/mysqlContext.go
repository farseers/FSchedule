package context

import (
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/trace"
)

// MysqlContextIns 初始化数据库上下文
var mysqlContextIns *mysqlContext

type mysqlContext struct {
	// 获取原生ORM框架（不使用TableSet或DomainSet）
	data.IInternalContext
	TaskGroup data.TableSet[model.TaskGroupPO] `data:"name=fschedule_task_group;migrate;"`
	Task      data.TableSet[model.TaskPO]      `data:"name=fschedule_task;migrate;"`
	TaskLog   data.TableSet[model.TaskLogPO]   `data:"name=fschedule_task_log;migrate;"`
}

// InitMysqlContext 初始化上下文
func InitMysqlContext() {
	mysqlContextIns = data.NewContext[mysqlContext]("default")
}

// 获取数据库上下文
func MysqlContextIns(cmt string) *mysqlContext {
	trace.SetComment(cmt)
	return mysqlContextIns
}
