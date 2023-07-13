package repository

import (
	"FSchedule/infrastructure/repository/context"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"strconv"
	"time"
)

type scheduleRepository struct {
}

func (receiver *scheduleRepository) ScheduleLock(name string, taskId int64) core.ILock {
	return context.RedisContextIns.LockNew("FSchedule_ScheduleLock:"+name+"_"+strconv.FormatInt(taskId, 10), strconv.FormatInt(fs.AppId, 10), 5*time.Second)
}

func (receiver *scheduleRepository) Election(fn func()) {
	go context.RedisContextIns.Election("FSchedule_Master", fn)
}

func (receiver *scheduleRepository) Schedule(name string, fn func()) {
	context.RedisContextIns.Election("FSchedule_Schedule:"+name, fn)
}

func (receiver *scheduleRepository) GetLeaderId() int64 {
	return context.RedisContextIns.GetLeaderId("FSchedule_Master")
}
