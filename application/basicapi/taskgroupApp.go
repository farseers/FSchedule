// @area /basicapi/
package basicapi

import (
	"FSchedule/application/basicapi/request"
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

// 任务组列表
// @get taskGroup/list
func TaskGroupList(name string, enable int, taskStatus enum.TaskStatus, clientId int64, pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[taskGroup.DomainObject] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return taskGroupRepository.ToListForPage(name, enable, taskStatus, clientId, pageSize, pageIndex)
}

// 任务组详情
// @get taskGroup/info-{taskGroupId}
func TaskGroupInfo(taskGroupId int64, taskGroupRepository taskGroup.Repository) taskGroup.DomainObject {
	return taskGroupRepository.ToEntity(taskGroupId)
}

// 任务组修改
// @get taskGroup/update
func TaskGroupUpdate(req request.TaskGroupUpdateRequest, taskGroupRepository taskGroup.Repository) {
	// 检查cron
	_, err := taskGroup.StandardParser.Parse(req.Cron)
	exception.ThrowWebExceptionError(403, err)
	// 判断任务组是否存在
	exception.ThrowWebExceptionBool(!taskGroupRepository.IsExists(req.Id), 403, "任务组不存在")
	// 更新
	taskGroupDO := mapper.Single[taskGroup.DomainObject](req)
	taskGroupDO.Update()
	taskGroupRepository.UpdateByEdit(taskGroupDO)
}

// 任务组删除
// @post taskGroup/delete
func TaskGroupDelete(taskGroupId int64, taskGroupRepository taskGroup.Repository) {
	taskGroupRepository.Delete(taskGroupId)
}

// 任务组数量
// @get taskGroup/count
func TaskGroupCount(taskGroupRepository taskGroup.Repository) int64 {
	return taskGroupRepository.GetTaskGroupCount()
}

// 任务组到期未运行数量
// @get taskGroup/unRunCount
func TaskGroupUnRunCount(taskGroupRepository taskGroup.Repository) int {
	return taskGroupRepository.GetUnRunCount()
}

// 任务组到期未运行任务组列表
// @get taskGroup/unRunList
func TaskGroupUnRunList(pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[taskGroup.DomainObject] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return taskGroupRepository.GetUnRunList(pageSize, pageIndex)
}

// 调度中或执行中的任务组
// @get taskGroup/schedulerWorkingList
func TaskGroupSchedulerList(pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[taskGroup.DomainObject] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return taskGroupRepository.ToSchedulerWorkingList(pageSize, pageIndex)
}
