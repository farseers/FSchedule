package repository

import (
	"FSchedule/infrastructure/repository/context"
	"strconv"
	"time"

	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/parse"
)

type scheduleRepository struct {
}

func (receiver *scheduleRepository) RegistryLock(clientId int64) core.ILock {
	return context.RedisContext("调度锁").LockNew("FSchedule_RegistryLock:"+parse.ToString(clientId), strconv.FormatInt(core.AppId, 10), 5*time.Second)
}

func (receiver *scheduleRepository) Election(fn func()) {
	context.RedisContext("选举").Election("FSchedule_Master", fn)
}

func (receiver *scheduleRepository) Schedule(taskGroupName string, fn func()) {
	context.RedisContext("任务组锁").Election("FSchedule_Schedule:"+taskGroupName, fn)
}

func (receiver *scheduleRepository) Monitor(fn func()) {
	context.RedisContext("监控任务组超时锁").Election("FSchedule_TaskGroupMonitor", fn)
}

func (receiver *scheduleRepository) GetLeaderId() int64 {
	return context.RedisContext("获取Master节点ID").GetLeaderId("FSchedule_Master")
}
