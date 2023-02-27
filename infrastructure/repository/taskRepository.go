package repository

import (
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type taskRepository struct {
	Task data.TableSet[model.TaskPO] `data:"name=task"`
}

func getCacheManager(name string) cache.ICacheManage[taskGroup.TaskEO] {
	lock.Lock()
	defer lock.Unlock()
	key := "FSchedule_Task:" + name
	if !container.IsRegister[cache.ICacheManage[taskGroup.TaskEO]](key) {
		cacheManage := redis.SetProfiles[taskGroup.TaskEO](key, "Id", 0, "default")
		cacheManage.SetItemSource(func(cacheId any) (taskGroup.TaskEO, bool) {
			flog.Warningf("我在读数据库:task %d", cacheId)
			repository := newManagerRepository()
			po := repository.Task.Where("Id = ?", cacheId).ToEntity()
			if po.Id > 0 {
				return mapper.Single[taskGroup.TaskEO](&po), true
			}
			return taskGroup.TaskEO{}, false
		})

		// 60秒同步一次任务到数据库
		syncTime := configure.GetInt("FSchedule.DataSyncTime")
		if syncTime > 0 {
			cacheManage.SetSyncSource(time.Duration(syncTime)*time.Second, func(do taskGroup.TaskEO) {
				// 保存成功后，已完成的任务，且最后运行时间大于1分钟的，移除列表
				// 最后运行时间超过1小时的移除。（如果有读取，还是会从数据库重新读的）
				if (do.IsFinish() && time.Now().Sub(do.RunAt).Seconds() >= float64(30)) ||
					(time.Now().Sub(do.RunAt).Hours() >= float64(1)) {
					flog.Debugf("同步数据库:共%d条，task:%d", cacheManage.Count(), do.Id)
					po := mapper.Single[model.TaskPO](&do)
					repository := newManagerRepository()
					result := repository.Task.UpdateOrInsert(po, "Id") == nil
					if result {
						cacheManage.Remove(po.Id)
						if cacheManage.ExistsItem(po.Id) {
							_ = flog.Error("数据删除失败")
						}
					}
				}
				time.Sleep(10 * time.Millisecond)
			})
		}
	}

	return container.Resolve[cache.ICacheManage[taskGroup.TaskEO]](key)
}

func (receiver *taskRepository) GetTask(name string, taskId int64) taskGroup.TaskEO {
	item, _ := getCacheManager(name).GetItem(taskId)
	return item
}

func (receiver *taskRepository) SaveTask(taskEO taskGroup.TaskEO) {
	getCacheManager(taskEO.Name).SaveItem(taskEO)
}

func (receiver *taskRepository) DeleteTask(name string) {
	receiver.Task.Where("name = ?", name).Delete()
	getCacheManager(name).Clear()
}

func (receiver *taskRepository) ToTaskSpeedList(name string) []int64 {
	lstPO := receiver.Task.Where("name = ? and status = ?", name, enum.Success).Desc("create_at").Select("RunSpeed").Limit(100).ToList()
	var lstSpeed []int64
	lstPO.Select(&lstSpeed, func(item model.TaskPO) any {
		return item.RunSpeed
	})
	return lstSpeed
}

func (receiver *taskRepository) ToFinishList(name string, top int) collections.List[taskGroup.TaskEO] {
	lstPO := receiver.Task.Where("name = ? and (status = ? or status = ?)", name, enum.Success, enum.Fail).Desc("create_at").Limit(top).ToList()
	return mapper.ToList[taskGroup.TaskEO](lstPO)
}

// ClearFinish 清除成功的任务记录（1天前）
func (receiver *taskRepository) ClearFinish(name string, taskId int) {
	receiver.Task.Where("task_group_name = ? and (status = ? or status = ?) and create_at < ? and Id < ?", name, enum.Success, enum.Fail, time.Now().Add(-24*time.Hour), taskId).Delete()
}
