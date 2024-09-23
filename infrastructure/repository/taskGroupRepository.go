package repository

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
	"FSchedule/domain/taskGroup"
	"FSchedule/infrastructure/repository/context"
	"FSchedule/infrastructure/repository/model"
	_ "embed"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
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
	taskGroupCache := redis.SetProfiles[taskGroup.DomainObject]("FSchedule_TaskGroup", "Name", "default")
	// 多级缓存
	taskGroupCache.SetListSource(func() collections.List[taskGroup.DomainObject] {
		list := context.MysqlContextIns("获取任务组列表").TaskGroup.ToList()
		return mapper.ToList[taskGroup.DomainObject](list)
	})

	taskGroupCache.SetItemSource(func(cacheId any) (taskGroup.DomainObject, bool) {
		po := context.MysqlContextIns("获取任务组").TaskGroup.Where("Name = ?", cacheId).ToEntity()
		if po.Name != "" {
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

func (receiver *taskGroupRepository) ToEntity(taskGroupName string) taskGroup.DomainObject {
	item, _ := receiver.CacheManage.GetItem(taskGroupName)
	return item
}

func (receiver *taskGroupRepository) Save(do taskGroup.DomainObject) {
	do.NeedSave = false
	receiver.CacheManage.SaveItem(do)

	// 发到所有节点上
	//_ = container.Resolve[core.IEvent]("TaskGroupUpdate").Publish(do)
}

func (receiver *taskGroupRepository) SaveAndTask(do taskGroup.DomainObject) {
	do.NeedSave = false
	receiver.CacheManage.SaveItem(do)
	receiver.SaveTask(do.Task)
}

func (receiver *taskGroupRepository) Sync() {
	lst := receiver.CacheManage.Get()

	for i := 0; i < lst.Count(); i++ {
		do := lst.Index(i)
		po := mapper.Single[model.TaskGroupPO](&do)

		if po.StartAt.Year() < 2000 {
			flog.Warningf("任务组：%s StartAt字段时间不正确 %s", do.Name, po.StartAt.String())
			po.StartAt = time.Now()
		}

		if po.ActivateAt.Year() < 2000 {
			flog.Warningf("任务组：%s ActivateAt字段时间不正确 %s", do.Name, po.ActivateAt.String())
			po.ActivateAt = time.Now()
		}

		if po.LastRunAt.Year() < 2000 {
			flog.Warningf("任务组：%s LastRunAt字段时间不正确 %s", do.Name, po.LastRunAt.String())
			po.LastRunAt = time.Now()
		}

		if po.NextAt.Year() < 2000 {
			flog.Warningf("任务组：%s NextAt字段时间不正确 %s", do.Name, po.NextAt.String())
			po.NextAt = time.Now()
		}
		_ = context.MysqlContextIns("更新任务组").TaskGroup.UpdateOrInsert(po, "name")

		// 同步任务
		receiver.taskRepository.syncTask(po.Name)
	}
}

func (receiver *taskGroupRepository) ToListForFops(taskGroupName string, enable int, taskStatus executeStatus.Enum, taskId int64, clientId string, pageSize int, pageIndex int) collections.List[taskGroup.DomainObject] {
	lst := receiver.CacheManage.Get()
	if taskGroupName != "" {
		lst = lst.Where(func(item taskGroup.DomainObject) bool {
			return strings.Contains(strings.ToLower(item.Name), strings.ToLower(taskGroupName)) || strings.Contains(strings.ToLower(item.Caption), strings.ToLower(taskGroupName))
		}).ToList()
	}
	if enable > -1 {
		lst = lst.Where(func(item taskGroup.DomainObject) bool {
			return item.IsEnable == parse.ToBool(enable)
		}).ToList()
	}
	if taskStatus > -1 {
		lst = lst.Where(func(item taskGroup.DomainObject) bool {
			return item.Task.ExecuteStatus == taskStatus
		}).ToList()
	}

	if clientId != "" {
		lst = lst.Where(func(item taskGroup.DomainObject) bool {
			return item.Task.Client.Id == clientId
		}).ToList()
	}
	if taskId > 0 {
		lst = lst.Where(func(item taskGroup.DomainObject) bool {
			return item.Task.Id == taskId
		}).ToList()
	}

	return lst
}

func (receiver *taskGroupRepository) IsExists(taskGroupName string) bool {
	return receiver.CacheManage.ExistsItem(taskGroupName)
}

func (receiver *taskGroupRepository) Delete(taskGroupName string) {
	// 删除任务
	(&taskRepository{}).DeleteTask(taskGroupName)
	// 删除日志
	(&TaskLogRepository{}).DeleteLog(taskGroupName)
	// 删除任务组
	_, _ = context.MysqlContextIns("删除任务组").TaskGroup.Where("name = ?", taskGroupName).Delete()
	// 删除缓存
	receiver.CacheManage.Remove(taskGroupName)
}

func (receiver *taskGroupRepository) GetTaskGroupCount() int64 {
	return int64(receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable
	}).Count())
}

func (receiver *taskGroupRepository) GetUnRunCount() int {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && (item.Task.ExecuteStatus == executeStatus.None && item.Task.ScheduleStatus == scheduleStatus.Scheduling) && item.NextAt.Before(dateTime.Now())
	}).Count()
}
