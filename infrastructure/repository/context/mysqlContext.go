package context

import (
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/data"
)

// MysqlContextIns 初始化数据库上下文
var MysqlContextIns *mysqlContext

type mysqlContext struct {
	TaskGroup data.TableSet[model.TaskGroupPO] `data:"name=fschedule_task_group"`
	Task      data.TableSet[model.TaskPO]      `data:"name=fschedule_task"`
	TaskLog   data.TableSet[model.TaskLogPO]   `data:"name=fschedule_task_log"`
}

// InitMysqlContext 初始化上下文
func InitMysqlContext() {
	MysqlContextIns = data.NewContext[mysqlContext]("default", true)
}
