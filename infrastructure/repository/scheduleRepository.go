package repository

import (
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/redis"
	"strconv"
	"time"
)

type scheduleRepository struct {
	*redis.Client
}

func (receiver *scheduleRepository) GetLock(name string, nextAt time.Time) core.ILock {
	return receiver.Lock.GetLocker("FSS_Scheduler:"+name+":"+strconv.FormatInt(nextAt.UnixMilli(), 10), 5*time.Second)
}
