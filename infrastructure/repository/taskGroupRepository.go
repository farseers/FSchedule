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
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
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
		var lst collections.List[taskGroup.DomainObject]
		list := context.MysqlContextIns.TaskGroup.ToList()
		list.MapToList(&lst)
		return lst
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
		do.ActivateAt = time.Now()
		do.LastRunAt = time.Now()
		do.NextAt = time.Now()
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
		_ = context.MysqlContextIns.TaskGroup.UpdateOrInsert(po, "id")

		// 同步任务
		receiver.taskRepository.syncTask(po.Id)
	}
}

func (receiver *taskGroupRepository) ToListByClientId(clientId int64) collections.List[taskGroup.DomainObject] {
	lst := receiver.CacheManage.Get()
	return lst.Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Client.Id == clientId && item.Task.StartAt.UnixMicro() < time.Now().UnixMicro()
	}).ToList()
}

func (receiver *taskGroupRepository) GetTaskGroupCount() int64 {
	return int64(receiver.CacheManage.Count())
}

func (receiver *taskGroupRepository) Delete(id int64) {
	context.MysqlContextIns.TaskGroup.Where("id = ?", id).Delete()
	receiver.CacheManage.Remove(id)
}

func (receiver *taskGroupRepository) ToUnRunCount() int {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Status == enum.None || item.Task.Status == enum.Scheduling || item.Task.CreateAt.UnixMicro() < time.Now().UnixMicro()
	}).Count()
}

func (receiver *taskGroupRepository) ToSchedulerWorkingList() collections.List[taskGroup.DomainObject] {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Status == enum.Scheduling || item.Task.Status == enum.Working
	}).ToList()
}

func (receiver *taskGroupRepository) GetTaskUnFinishList(jobsNames []string, top int) collections.List[taskGroup.DomainObject] {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && collections.NewList(jobsNames...).Contains(item.Name) && item.Task.Status != enum.Success && item.Task.Status != enum.Fail
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.NextAt.UnixMicro()
	}).Take(top).ToList()
}

// SaveToDb 保存到数据库
func (receiver *taskGroupRepository) SaveToDb(do taskGroup.DomainObject) {
	po := mapper.Single[model.TaskGroupPO](&do)
	context.MysqlContextIns.TaskGroup.Where("id = ?", do.Id).Update(po)
}

// ToIdList 从数据库中读取数据
func (receiver *taskGroupRepository) ToIdList() []int64 {
	lst := context.MysqlContextIns.TaskGroup.Select("id").ToList()
	var lstId []int64
	lst.Select(&lstId, func(item model.TaskGroupPO) any {
		return item.Id
	})
	return lstId
}

func (receiver *taskGroupRepository) GetEnableTaskList(status enum.TaskStatus, pageSize int, pageIndex int) collections.PageList[taskGroup.TaskEO] {
	lstTaskGroup := receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable
	}).ToList()

	if status != enum.None {
		lstTaskGroup = lstTaskGroup.Where(func(item taskGroup.DomainObject) bool {
			return item.Task.Status == status
		}).ToList()
	}

	lstTaskGroup = lstTaskGroup.OrderBy(func(item taskGroup.DomainObject) any {
		return item.Name
	}).ToList()

	var lst collections.List[taskGroup.TaskEO]
	lstTaskGroup.Select(&lst, func(item taskGroup.DomainObject) any {
		return item.Task
	})
	return lst.ToPageList(pageSize, pageIndex)
}
