package repository

import (
	"FSchedule/domain/taskGroup"
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
)

type taskRepository struct {
	Task data.TableSet[model.TaskPO] `data:"name=task"`
}

func getCacheManager(name string) cache.ICacheManage[taskGroup.TaskEO] {
	key := "FSchedule_Task:" + name
	if !container.IsRegister[taskGroup.TaskEO](key) {
		profiles := redis.SetProfiles[taskGroup.TaskEO](key, "Id", 0, "default")
		profiles.SetItemSource(func(cacheId any) (taskGroup.TaskEO, bool) {
			repository := data.NewContext[taskRepository]("default")
			po := repository.Task.Where("Id = ?", cacheId).ToEntity()
			if po.Id > 0 {
				return mapper.Single[taskGroup.TaskEO](&po), true
			}
			return taskGroup.TaskEO{}, false
		})
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

func (receiver *taskRepository) AddTask(taskDO taskGroup.TaskEO) {
	po := mapper.Single[model.TaskPO](&taskDO)
	receiver.Task.Insert(&po)
}
