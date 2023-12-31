package repository

import (
	"FSchedule/infrastructure/repository/context"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/parse"
	"strconv"
	"time"
)

type scheduleRepository struct {
}

func (receiver *scheduleRepository) ScheduleLock(taskGroupId int64, taskId int64) core.ILock {
	return context.RedisContextIns.LockNew("FSchedule_ScheduleLock:"+parse.ToString(taskGroupId)+"_"+parse.ToString(taskId), strconv.FormatInt(core.AppId, 10), 5*time.Second)
}

func (receiver *scheduleRepository) Election(fn func()) {
	go context.RedisContextIns.Election("FSchedule_Master", fn)
}

func (receiver *scheduleRepository) Schedule(taskGroupId int64, fn func()) {
	context.RedisContextIns.Election("FSchedule_Schedule:"+parse.ToString(taskGroupId), fn)
}

func (receiver *scheduleRepository) GetLeaderId() int64 {
	return context.RedisContextIns.GetLeaderId("FSchedule_Master")
}
