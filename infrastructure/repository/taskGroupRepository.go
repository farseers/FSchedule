package repository

import (
	"FSchedule/domain/taskGroup"
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
	"time"
)

type taskGroupRepository struct {
	Redis                   redis.IClient `inject:"default"`
	TaskGroupUpdateEventBus core.IEvent   `inject:"ClientUpdate"`
	*managerRepository
	*taskRepository
}

func registerTaskGroupRepository() {
	cacheManage := redis.SetProfiles[taskGroup.DomainObject]("FSchedule_TaskGroup", "Name", 0, "default")
	cacheManage.EnableItemNullToLoadAll()
	// 多级缓存
	cacheManage.SetListSource(func() collections.List[taskGroup.DomainObject] {
		var lst collections.List[taskGroup.DomainObject]
		repository := newManagerRepository()
		list := repository.TaskGroup.ToList()
		list.MapToList(&lst)
		return lst
	})

	cacheManage.SetItemSource(func(cacheId any) (taskGroup.DomainObject, bool) {
		repository := newManagerRepository()
		po := repository.TaskGroup.Where("Name = ?", cacheId).ToEntity()
		if po.Name != "" {
			return mapper.Single[taskGroup.DomainObject](&po), true
		}
		return taskGroup.DomainObject{}, false
	})

	// 60秒同步一次任务组到数据库
	cacheManage.SetSyncSource(time.Duration(configure.GetInt("FSchedule.DataSyncTime"))*time.Second, func(do taskGroup.DomainObject) {
		po := mapper.Single[model.TaskGroupPO](&do)
		repository := newManagerRepository()
		_ = repository.TaskGroup.UpdateOrInsert(po, "Name")
	})

	// 注册仓储
	container.Register(func() taskGroup.Repository {
		repository := container.ResolveIns(&taskGroupRepository{})
		repository.managerRepository = newManagerRepository()
		repository.taskRepository = &taskRepository{
			Task: repository.managerRepository.Task,
		}
		return repository
	})
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
