package repository

import (
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"FSchedule/infrastructure/repository/context"
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
	"strings"
	"time"
)

type taskGroupRepository struct {
	CacheManage cache.ICacheManage[taskGroup.DomainObject] `inject:"FSchedule_TaskGroup"`
	*taskRepository
}

func registerTaskGroupRepository() {
	taskGroupCache := redis.SetProfiles[taskGroup.DomainObject]("FSchedule_TaskGroup", "Id", "default")
	// 多级缓存
	taskGroupCache.SetListSource(func() collections.List[taskGroup.DomainObject] {
		list := context.MysqlContextIns.TaskGroup.ToList()
		return mapper.ToList[taskGroup.DomainObject](list)
	})

	taskGroupCache.SetItemSource(func(cacheId any) (taskGroup.DomainObject, bool) {
		po := context.MysqlContextIns.TaskGroup.Where("Id = ?", cacheId).ToEntity()
		if po.Name != "" && po.Id > 0 {
			return mapper.Single[taskGroup.DomainObject](&po), true
		}
		return taskGroup.DomainObject{}, false
	})

	// 注册仓储
	container.Register(func() taskGroup.Repository {
		return &taskGroupRepository{taskRepository: &taskRepository{}}
	})
}

func (receiver *taskGroupRepository) ToList() collections.List[taskGroup.DomainObject] {
	return receiver.CacheManage.Get()
}

func (receiver *taskGroupRepository) ToListByName(name string) collections.List[taskGroup.DomainObject] {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.Name == name
	}).ToList()
}

func (receiver *taskGroupRepository) ToEntity(id int64) taskGroup.DomainObject {
	item, _ := receiver.CacheManage.GetItem(id)
	return item
}

func (receiver *taskGroupRepository) Save(do taskGroup.DomainObject) {
	do.NeedSave = false
	// 说明是新注册的任务
	if do.Id == 0 {
		do.ActivateAt = dateTime.Now()
		do.LastRunAt = dateTime.Now()
		do.NextAt = dateTime.Now()
		po := mapper.Single[model.TaskGroupPO](&do)
		_ = context.MysqlContextIns.TaskGroup.Insert(&po)
		do.Id = po.Id
		do.Task.TaskGroupId = po.Id
	}
	receiver.CacheManage.SaveItem(do)

	// 发到所有节点上
	_ = container.Resolve[core.IEvent]("TaskGroupUpdate").Publish(do)
}

func (receiver *taskGroupRepository) SaveAndTask(do taskGroup.DomainObject) {
	do.NeedSave = false
	receiver.Save(do)
	receiver.SaveTask(do.Task)
}

func (receiver *taskGroupRepository) Sync() {
	lst := receiver.CacheManage.Get()
	for i := 0; i < lst.Count(); i++ {
		do := lst.Index(i)
		po := mapper.Single[model.TaskGroupPO](&do)

		if po.StartAt.Year() < 2000 {
			po.StartAt = time.Now()
		}

		if po.ActivateAt.Year() < 2000 {
			po.ActivateAt = time.Now()
		}

		if po.LastRunAt.Year() < 2000 {
			po.LastRunAt = time.Now()
		}

		if po.NextAt.Year() < 2000 {
			po.NextAt = time.Now()
		}
		_ = context.MysqlContextIns.TaskGroup.UpdateOrInsert(po, "id")

		// 同步任务
		receiver.taskRepository.syncTask(po.Id)
	}
}

func (receiver *taskGroupRepository) ToListForPage(name string, enable int, taskStatus enum.TaskStatus, clientId int64, pageSize int, pageIndex int) collections.PageList[taskGroup.DomainObject] {
	lst := receiver.CacheManage.Get()
	if name != "" {
		lst = lst.Where(func(item taskGroup.DomainObject) bool {
			return strings.Contains(strings.ToLower(item.Name), strings.ToLower(name))
		}).ToList()
	}
	if enable > -1 {
		lst = lst.Where(func(item taskGroup.DomainObject) bool {
			return item.IsEnable == parse.ToBool(enable)
		}).ToList()
	}
	if taskStatus > -1 {
		lst = lst.Where(func(item taskGroup.DomainObject) bool {
			return item.Task.Status == taskStatus
		}).ToList()
	}

	if clientId > 0 {
		lst = lst.Where(func(item taskGroup.DomainObject) bool {
			return item.Task.Client.Id == clientId
		}).ToList()
	}

	// 排序
	return lst.OrderBy(func(item taskGroup.DomainObject) any {
		return item.Name
	}).ToPageList(pageSize, pageIndex)
}

func (receiver *taskGroupRepository) IsExists(taskGroupId int64) bool {
	return receiver.CacheManage.ExistsItem(taskGroupId)
}

func (receiver *taskGroupRepository) UpdateByEdit(do taskGroup.DomainObject) {
	if item, exists := receiver.CacheManage.GetItem(do.Id); exists {
		item.Name = do.Name
		item.Ver = do.Ver
		item.Caption = do.Caption
		item.Data = do.Data
		item.StartAt = do.StartAt
		item.NextAt = do.NextAt
		item.Cron = do.Cron
		item.IsEnable = do.IsEnable
		receiver.Save(do)
	}
}

func (receiver *taskGroupRepository) Delete(taskGroupId int64) {
	(&taskRepository{}).DeleteTask(taskGroupId)
	_, _ = context.MysqlContextIns.TaskGroup.Where("id = ?", taskGroupId).Delete()
	receiver.CacheManage.Remove(taskGroupId)
}

func (receiver *taskGroupRepository) GetTaskGroupCount() int64 {
	return int64(receiver.CacheManage.Count())
}

func (receiver *taskGroupRepository) GetUnRunCount() int {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && (item.Task.Status == enum.None || item.Task.Status == enum.Scheduling) && item.NextAt.Before(dateTime.Now())
	}).Count()
}

func (receiver *taskGroupRepository) GetUnRunList(pageSize int, pageIndex int) collections.PageList[taskGroup.DomainObject] {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && (item.Task.Status == enum.None || item.Task.Status == enum.Scheduling) && item.NextAt.Before(dateTime.Now())
	}).ToPageList(pageSize, pageIndex)
}

func (receiver *taskGroupRepository) ToSchedulerWorkingList(pageSize int, pageIndex int) collections.PageList[taskGroup.DomainObject] {
	lst := receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Status == enum.Scheduling || item.Task.Status == enum.Working
	}).ToList()
	return lst.ToPageList(pageSize, pageIndex)
}
