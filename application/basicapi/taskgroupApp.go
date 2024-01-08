// @area /basicapi/taskGroup/
package basicapi

import (
	"FSchedule/application/basicapi/request"
	"FSchedule/application/basicapi/response"
	"FSchedule/domain"
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

// 任务组列表
// @get list
func TaskGroupList(taskGroupName string, enable int, taskStatus enum.TaskStatus, clientId int64, pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[response.TaskGroupResponse] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	lst := taskGroupRepository.ToListForPage(taskGroupName, enable, taskStatus, clientId, pageSize, pageIndex)
	lstTaskGroupResponse := mapper.ToPageList[response.TaskGroupResponse](lst)
	// 获取每个任务组当前注册的客户端
	lstTaskGroupResponse.List.Foreach(func(item *response.TaskGroupResponse) {
		item.Clients = collections.NewList[response.ClientResponse]()
		for _, c := range domain.GetClientList(item.Name).ToArray() {
			item.Clients.Add(mapper.Single[response.ClientResponse](*c))
		}
	})
	return lstTaskGroupResponse
}

// 任务组详情
// @get info-{taskGroupId}
func TaskGroupInfo(taskGroupName string, taskGroupRepository taskGroup.Repository) response.TaskGroupResponse {
	// 判断任务组是否存在
	exception.ThrowWebExceptionBool(!taskGroupRepository.IsExists(taskGroupName), 403, "任务组不存在")

	info := taskGroupRepository.ToEntity(taskGroupName)
	item := mapper.Single[response.TaskGroupResponse](info)
	item.Clients = collections.NewList[response.ClientResponse]()
	for _, c := range domain.GetClientList(item.Name).ToArray() {
		item.Clients.Add(mapper.Single[response.ClientResponse](*c))
	}
	return item
}

// 任务组修改
// @post update
func TaskGroupUpdate(req request.TaskGroupUpdateRequest, taskGroupRepository taskGroup.Repository) {
	// 检查cron
	_, err := taskGroup.StandardParser.Parse(req.Cron)
	exception.ThrowWebExceptionError(403, err)
	// 判断任务组是否存在
	exception.ThrowWebExceptionBool(!taskGroupRepository.IsExists(req.Name), 403, "任务组不存在")
	// 更新
	taskGroupDO := mapper.Single[taskGroup.DomainObject](req)
	taskGroupDO.Update()
	taskGroupRepository.UpdateByEdit(taskGroupDO)
}

// 任务组删除
// @post delete
func TaskGroupDelete(taskGroupName string, taskGroupRepository taskGroup.Repository) {
	// 判断任务组是否存在
	exception.ThrowWebExceptionBool(!taskGroupRepository.IsExists(taskGroupName), 403, "任务组不存在")

	taskGroupRepository.Delete(taskGroupName)
}

// 任务组数量
// @get count
func TaskGroupCount(taskGroupRepository taskGroup.Repository) int64 {
	return taskGroupRepository.GetTaskGroupCount()
}

// 任务组到期未运行数量
// @get unRunCount
func TaskGroupUnRunCount(taskGroupRepository taskGroup.Repository) int {
	return taskGroupRepository.GetUnRunCount()
}

// 任务组到期未运行任务组列表
// @get unRunList
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
// @get schedulerWorkingList
func TaskGroupSchedulerList(pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[taskGroup.DomainObject] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return taskGroupRepository.ToSchedulerWorkingList(pageSize, pageIndex)
}

// 设置任务组状态
// @post setEnable
func SetEnable(taskGroupName string, enable bool, taskGroupRepository taskGroup.Repository) {
	taskGroupDO := taskGroupRepository.ToEntity(taskGroupName)
	taskGroupDO.SetEnable(enable)
	taskGroupRepository.Save(taskGroupDO)
}
