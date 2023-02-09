package repository

import (
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"FSchedule/infrastructure/repository/model"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/mapper"
	"time"
)

type managerRepository struct {
	Task        data.TableSet[model.TaskPO]                `data:"name=task"`
	TaskGroup   data.TableSet[model.TaskGroupPO]           `data:"name=task_group"`
	CacheManage cache.ICacheManage[taskGroup.DomainObject] `inject:"FSS_TaskGroup"`
}

func (receiver *managerRepository) ToListByClientId(clientId int64) collections.List[taskGroup.DomainObject] {
	lst := receiver.CacheManage.Get()
	return lst.Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Client.Id == clientId && item.Task.StartAt.UnixMicro() < time.Now().UnixMicro()
	}).ToList()
}

func (receiver *managerRepository) GetTaskGroupCount() int64 {
	return int64(receiver.CacheManage.Count())
}

func (receiver *managerRepository) Add(do *taskGroup.DomainObject) {
	po := mapper.Single[model.TaskGroupPO](do)
	po.ActivateAt = time.Now()
	po.LastRunAt = time.Now()
	po.NextAt = time.Now()
	receiver.TaskGroup.Insert(&po)
	receiver.CacheManage.SaveItem(*do)
}

func (receiver *managerRepository) Delete(name string) {
	receiver.TaskGroup.Where("name = ?", name).Delete()
	receiver.CacheManage.Remove(name)
}

func (receiver *managerRepository) ToUnRunCount() int {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Status == enum.None || item.Task.Status == enum.Scheduling || item.Task.CreateAt.UnixMicro() < time.Now().UnixMicro()
	}).Count()
}

func (receiver *managerRepository) ToSchedulerWorkingList() collections.List[taskGroup.DomainObject] {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.Task.Status == enum.Scheduling || item.Task.Status == enum.Working
	}).ToList()
}

func (receiver *managerRepository) GetTaskUnFinishList(jobsNames []string, top int) collections.List[taskGroup.DomainObject] {
	return receiver.CacheManage.Get().Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && collections.NewList(jobsNames...).Contains(item.Name) && item.Task.Status != enum.Success && item.Task.Status != enum.Fail
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.NextAt.UnixMicro()
	}).Take(top).ToList()
}

// SaveToDb 保存到数据库
func (receiver *managerRepository) SaveToDb(do taskGroup.DomainObject) {
	po := mapper.Single[model.TaskGroupPO](&do)
	receiver.TaskGroup.Where("name = ?", do.Name).Update(po)
}

// ToIdList 从数据库中读取数据
func (receiver *managerRepository) ToIdList() []string {
	lst := receiver.TaskGroup.Select("name").ToList()
	var lstName []string
	lst.Select(&lstName, func(item model.TaskGroupPO) any {
		return item.Name
	})
	return lstName
}

func (receiver *managerRepository) GetEnableTaskList(status enum.TaskStatus, pageSize int, pageIndex int) collections.PageList[taskGroup.TaskEO] {
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

func (receiver *managerRepository) ToTaskSpeedList(name string) []int64 {
	lstPO := receiver.Task.Where("name = ? and status = ?", name, enum.Success).Desc("create_at").Select("RunSpeed").Limit(100).ToList()
	var lstSpeed []int64
	lstPO.Select(&lstSpeed, func(item model.TaskPO) any {
		return item.RunSpeed
	})
	return lstSpeed
}

func (receiver *managerRepository) TodayFailCount() int64 {
	return receiver.Task.Where("status = ? and create_at >= ?", enum.Fail, dateTime.Now().Date().ToTime()).Count()
}

// ClearFinish 清除成功的任务记录（1天前）
func (receiver *managerRepository) ClearFinish(name string, taskId int) {
	receiver.Task.Where("name = ? and (status = ? or status = ?) and create_at < ? and Id < ?", name, enum.Success, enum.Fail, time.Now().Add(-24*time.Hour), taskId).Delete()
}

func (receiver *managerRepository) ToListByGroupId(name string, pageSize int, pageIndex int) collections.PageList[taskGroup.TaskEO] {
	page := receiver.Task.Where("name = ?", name).Desc("create_at").ToPageList(pageSize, pageIndex)
	return receiver.toPageListTaskEO(page)
}
func (receiver *managerRepository) ToFinishPageList(pageSize int, pageIndex int) collections.PageList[taskGroup.TaskEO] {
	page := receiver.Task.Where("(status = ? or status = ?) and (create_at >= ?)", enum.Fail, enum.Success, time.Now().Add(-24*time.Hour)).
		Desc("run_at").ToPageList(pageSize, pageIndex)
	return receiver.toPageListTaskEO(page)
}

func (receiver *managerRepository) ToFinishList(name string, top int) collections.List[taskGroup.TaskEO] {
	lstPO := receiver.Task.Where("name = ? and (status = ? or status = ?)", name, enum.Success, enum.Fail).Desc("create_at").Limit(top).ToList()
	return mapper.ToList[taskGroup.TaskEO](lstPO)
}

func (receiver *managerRepository) toPageListTaskEO(page collections.PageList[model.TaskPO]) collections.PageList[taskGroup.TaskEO] {
	lst := mapper.ToList[taskGroup.TaskEO](page.List)
	return collections.NewPageList[taskGroup.TaskEO](lst, page.RecordCount)
}
