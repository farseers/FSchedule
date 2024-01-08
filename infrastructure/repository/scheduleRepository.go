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

func (receiver *scheduleRepository) ScheduleLock(taskGroupName string, taskId int64) core.ILock {
	return context.RedisContextIns.LockNew("FSchedule_ScheduleLock:"+taskGroupName+"_"+parse.ToString(taskId), strconv.FormatInt(core.AppId, 10), 5*time.Second)
}

func (receiver *scheduleRepository) Election(fn func()) {
	go context.RedisContextIns.Election("FSchedule_Master", fn)
}

func (receiver *scheduleRepository) Schedule(taskGroupName string, fn func()) {
	context.RedisContextIns.Election("FSchedule_Schedule:"+taskGroupName, fn)
}

func (receiver *scheduleRepository) GetLeaderId() int64 {
	return context.RedisContextIns.GetLeaderId("FSchedule_Master")
}
