package repository

import (
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"FSchedule/infrastructure/repository/context"
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type taskRepository struct {
}

func getCacheManager(taskGroupId int64) cache.ICacheManage[taskGroup.TaskEO] {
	key := "FSchedule_Task:" + parse.ToString(taskGroupId)
	if !container.IsRegister[cache.ICacheManage[taskGroup.TaskEO]](key) {
		lock.Lock()
		defer lock.Unlock()
		if !container.IsRegister[cache.ICacheManage[taskGroup.TaskEO]](key) {
			cacheManage := redis.SetProfiles[taskGroup.TaskEO](key, "Id", "default")
			cacheManage.SetItemSource(func(cacheId any) (taskGroup.TaskEO, bool) {
				po := context.MysqlContextIns.Task.Where("Id = ?", cacheId).ToEntity()
				if po.Id > 0 {
					return mapper.Single[taskGroup.TaskEO](&po), true
				}
				return taskGroup.TaskEO{}, false
			})
		}
	}

	return container.Resolve[cache.ICacheManage[taskGroup.TaskEO]](key)
}

func (receiver *taskRepository) GetTask(taskGroupId int64, taskId int64) taskGroup.TaskEO {
	item, _ := getCacheManager(taskGroupId).GetItem(taskId)
	return item
}

func (receiver *taskRepository) SaveTask(taskEO taskGroup.TaskEO) {
	getCacheManager(taskEO.TaskGroupId).SaveItem(taskEO)
}

func (receiver *taskRepository) syncTask(taskGroupId int64) {
	cacheManager := getCacheManager(taskGroupId)
	lst := cacheManager.Get()
	for i := 0; i < lst.Count(); i++ {
		do := lst.Index(i)
		flog.Infof("同步数据库:%d/%d，task:%d", i+1, lst.Count(), do.Id)
		// 保存成功后，已完成的任务，且最后运行时间大于1分钟的，移除列表
		// 最后运行时间超过1小时的移除。（如果有读取，还是会从数据库重新读的）
		if (do.IsFinish() && time.Now().Sub(do.RunAt).Seconds() >= float64(30)) ||
			(time.Now().Sub(do.RunAt).Hours() >= float64(1)) {
			po := mapper.Single[model.TaskPO](&do)
			if context.MysqlContextIns.Task.UpdateOrInsert(po, "Id") == nil {
				cacheManager.Remove(po.Id)
			}
		}
	}
}

func (receiver *taskRepository) DeleteTask(taskGroupId int64) {
	context.MysqlContextIns.Task.Where("task_group_id = ?", taskGroupId).Delete()
	getCacheManager(taskGroupId).Clear()
}

func (receiver *taskRepository) ToTaskListByGroupId(taskGroupId int64, taskStatus enum.TaskStatus, pageSize int, pageIndex int) collections.PageList[taskGroup.TaskEO] {
	ts := context.MysqlContextIns.Task.Desc("create_at")
	if taskGroupId > 0 {
		ts = ts.Where("task_group_id = ?", taskGroupId)
	}
	if taskStatus > -1 {
		ts = ts.Where("status = ?", taskStatus)
	}
	lstPO := ts.ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[taskGroup.TaskEO](lstPO)
}

func (receiver *taskRepository) TodayFailCount() int64 {
	now := dateTime.Now()
	return context.MysqlContextIns.Task.Where("status = ? and create_at >= ?", enum.Fail, now.Date()).Count()
}

func (receiver *taskRepository) ToTaskSpeedList(taskGroupId int64) []int64 {
	lstPO := context.MysqlContextIns.Task.Where("task_group_id = ? and status = ?", taskGroupId, enum.Success).Desc("create_at").Select("RunSpeed").Limit(100).ToList()
	var lstSpeed []int64
	lstPO.Select(&lstSpeed, func(item model.TaskPO) any {
		return item.RunSpeed
	})
	return lstSpeed
}

// TaskClearFinish 清除成功的任务记录（1天前）
func (receiver *taskRepository) TaskClearFinish(taskGroupId int64, taskId int) {
	context.MysqlContextIns.Task.Where("task_group_id = ? and (status = ? or status = ?) and create_at < ? and Id < ?", taskGroupId, enum.Success, enum.Fail, time.Now().Add(-24*time.Hour), taskId).Delete()
}

func (receiver *taskRepository) ToTaskFinishList(taskGroupId int64, top int) collections.List[taskGroup.TaskEO] {
	lstPO := context.MysqlContextIns.Task.Where("task_group_id = ? and (status = ? or status = ?)", taskGroupId, enum.Success, enum.Fail).Desc("create_at").Limit(top).ToList()
	return mapper.ToList[taskGroup.TaskEO](lstPO)

}
func (receiver *taskRepository) ToTaskFinishPageList(pageSize int, pageIndex int) collections.PageList[taskGroup.TaskEO] {
	page := context.MysqlContextIns.Task.Where("(status = ? or status = ?) and (create_at >= ?)", enum.Fail, enum.Success, time.Now().Add(-24*time.Hour)).
		Desc("run_at").ToPageList(pageSize, pageIndex)
	return receiver.toTaskPageListTaskEO(page)
}

func (receiver *taskRepository) toTaskPageListTaskEO(page collections.PageList[model.TaskPO]) collections.PageList[taskGroup.TaskEO] {
	lst := mapper.ToList[taskGroup.TaskEO](page.List)
	return collections.NewPageList[taskGroup.TaskEO](lst, page.RecordCount)
}
