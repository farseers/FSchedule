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
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
	"time"
)

type taskGroupRepository struct {
	TaskGroup               data.TableSet[model.TaskGroupPO]           `data:"name=test_task_group"`
	Redis                   redis.IClient                              `inject:"default"`
	TaskGroupUpdateEventBus core.IEvent                                `inject:"TaskGroupUpdate"`
	CacheManage             cache.ICacheManage[taskGroup.DomainObject] `inject:"FSchedule_TaskGroup"`
	*taskRepository
}

func registerTaskGroupRepository() {
	repository := data.NewContext[taskGroupRepository]("default", true)
	repository.taskRepository = data.NewContext[taskRepository]("default", true)

	repository.CacheManage = redis.SetProfiles[taskGroup.DomainObject]("FSchedule_TaskGroup", "Name", 0, "default")
	// 多级缓存
	repository.CacheManage.SetListSource(func() collections.List[taskGroup.DomainObject] {
		var lst collections.List[taskGroup.DomainObject]
		list := repository.TaskGroup.ToList()
		list.MapToList(&lst)
		return lst
	})

	repository.CacheManage.SetItemSource(func(cacheId any) (taskGroup.DomainObject, bool) {
		po := repository.TaskGroup.Where("Name = ?", cacheId).ToEntity()
		if po.Name != "" {
			return mapper.Single[taskGroup.DomainObject](&po), true
		}
		return taskGroup.DomainObject{}, false
	})

	// 60秒同步一次任务组到数据库
	syncTime := configure.GetInt("FSchedule.DataSyncTime")
	if syncTime > 0 {
		repository.CacheManage.SetSyncSource(time.Duration(syncTime)*time.Second, func(do taskGroup.DomainObject) {
			po := mapper.Single[model.TaskGroupPO](&do)
			_ = repository.TaskGroup.UpdateOrInsert(po, "Name")
		})
	}

	*repository = *container.ResolveIns(repository)

	// 注册仓储
	container.RegisterInstance[taskGroup.Repository](repository)
}

func (receiver *taskGroupRepository) Add(do *taskGroup.DomainObject) {
	po := mapper.Single[model.TaskGroupPO](do)
	po.ActivateAt = time.Now()
	po.LastRunAt = time.Now()
	po.NextAt = time.Now()
	_ = receiver.TaskGroup.Insert(&po)
	receiver.CacheManage.SaveItem(*do)
}

func (receiver *taskGroupRepository) ToList() collections.List[taskGroup.DomainObject] {
	return receiver.CacheManage.Get()
}

func (receiver *taskGroupRepository) ToEntity(name string) taskGroup.DomainObject {
	item, _ := receiver.CacheManage.GetItem(name)
	return item
}

func (receiver *taskGroupRepository) Save(do taskGroup.DomainObject) {
	do.NeedSave = false
	receiver.CacheManage.SaveItem(do)

	// 发到所有节点上
	_ = receiver.TaskGroupUpdateEventBus.Publish(do)
}

func (receiver *taskGroupRepository) SaveAndTask(do taskGroup.DomainObject) {
	do.NeedSave = false
	receiver.Save(do)
	receiver.SaveTask(do.Task)
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

func (receiver *taskGroupRepository) Delete(name string) {
	receiver.TaskGroup.Where("name = ?", name).Delete()
	receiver.CacheManage.Remove(name)
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
	receiver.TaskGroup.Where("name = ?", do.Name).Update(po)
}

// ToIdList 从数据库中读取数据
func (receiver *taskGroupRepository) ToIdList() []string {
	lst := receiver.TaskGroup.Select("name").ToList()
	var lstName []string
	lst.Select(&lstName, func(item model.TaskGroupPO) any {
		return item.Name
	})
	return lstName
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
