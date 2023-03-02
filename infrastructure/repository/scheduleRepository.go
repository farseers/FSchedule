package repository

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/redis"
	"strconv"
	"time"
)

type scheduleRepository struct {
	redis.IClient `inject:"default"`
}

func (receiver *scheduleRepository) ScheduleLock(name string, taskId int64) core.ILock {
	return receiver.LockNew("FSchedule_ScheduleLock:"+name+"_"+strconv.FormatInt(taskId, 10), strconv.FormatInt(fs.AppId, 10), 5*time.Second)
}

func (receiver *scheduleRepository) Election(fn func()) {
	go receiver.IClient.Election("FSchedule_Master", fn)
}

func (receiver *scheduleRepository) Schedule(name string, fn func()) {
	receiver.IClient.Election("FSchedule_Schedule:"+name, fn)
}

func (receiver *scheduleRepository) GetLeaderId() int64 {
	return receiver.IClient.GetLeaderId("FSchedule_Master")
}
