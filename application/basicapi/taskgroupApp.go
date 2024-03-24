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
func TaskGroupList(clientName, taskGroupName string, enable int, taskStatus enum.TaskStatus, taskId, clientId int64, pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[response.TaskGroupResponse] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	lst := taskGroupRepository.ToListForPage(clientName, taskGroupName, enable, taskStatus, taskId, clientId, pageSize, pageIndex)
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
// @get info-{taskGroupName}
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
	// 确认cron格式是否正确
	_, err := taskGroup.StandardParser.Parse(req.Cron)
	if err != nil {
		exception.ThrowWebExceptionf(403, "任务组:%s %s，Cron格式[%s]错误:%s", req.Name, req.Caption, req.Cron, err.Error())
	}

	// 判断任务组是否存在
	taskGroupDO := taskGroupRepository.ToEntity(req.Name)
	exception.ThrowWebExceptionBool(taskGroupDO.IsNil(), 403, "任务组不存在")

	err = mapper.Auto(req, &taskGroupDO)
	exception.ThrowWebExceptionError(403, err)

	// 更新
	taskGroupDO.Update()
	taskGroupRepository.Save(taskGroupDO)
}

// 任务组删除
// @post delete
func TaskGroupDelete(taskGroupName string, taskGroupRepository taskGroup.Repository) {
	// 判断任务组是否存在
	exception.ThrowWebExceptionBool(!taskGroupRepository.IsExists(taskGroupName), 403, "任务组不存在")

	taskGroupRepository.Delete(taskGroupName)
}

// 设置任务组状态
// @post setEnable
func SetEnable(taskGroupName string, enable bool, taskGroupRepository taskGroup.Repository) {
	taskGroupDO := taskGroupRepository.ToEntity(taskGroupName)
	taskGroupDO.SetEnable(enable)
	taskGroupRepository.Save(taskGroupDO)
}

/*
// 任务组到期未运行任务组列表
// @get unRunList
func TaskGroupUnRunList(pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[response.TaskGroupResponse] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}

	lst := taskGroupRepository.GetUnRunList(pageSize, pageIndex)
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

// 调度中或执行中的任务组
// @get schedulerWorkingList
func TaskGroupSchedulerList(pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[response.TaskGroupResponse] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	lst := taskGroupRepository.ToSchedulerWorkingList(pageSize, pageIndex)
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
*/
