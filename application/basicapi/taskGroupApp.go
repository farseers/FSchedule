// @area /basicapi/taskGroup/
package basicapi

import (
	"FSchedule/application/basicapi/request"
	"FSchedule/application/basicapi/response"
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

// 任务组列表
// @get list
func TaskGroupList(clientName, taskGroupName string, enable int, taskStatus executeStatus.Enum, taskId, clientId int64, pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository, clientRepository client.Repository) collections.PageList[response.TaskGroupResponse] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}

	lst := taskGroupRepository.ToListForFops(taskGroupName, enable, taskStatus, taskId, clientId, pageSize, pageIndex)
	lstTaskGroupResponse := mapper.ToList[response.TaskGroupResponse](lst)

	// 匹配每个任务组的客户端
	lstClient := clientRepository.ToList()
	lstTaskGroupResponse.Foreach(func(item *response.TaskGroupResponse) {
		item.Clients = mapper.ToList[response.ClientResponse](getTaskGroupClientList(item.Name, lstClient))
	})

	// 筛选客户端
	if clientName != "" {
		lstTaskGroupResponse = lstTaskGroupResponse.Where(func(item response.TaskGroupResponse) bool {
			return item.Clients.Where(func(client response.ClientResponse) bool {
				return client.Id == clientName || client.Name == clientName
			}).Any()
		}).ToList()
	}

	// 排序
	return lstTaskGroupResponse.OrderBy(func(item response.TaskGroupResponse) any {
		return item.Name
	}).ToPageList(pageSize, pageIndex)
}

// 任务组详情
// @get info-{taskGroupName}
func TaskGroupInfo(taskGroupName string, taskGroupRepository taskGroup.Repository, clientRepository client.Repository) response.TaskGroupResponse {
	// 判断任务组是否存在
	exception.ThrowWebExceptionBool(!taskGroupRepository.IsExists(taskGroupName), 403, "任务组不存在")

	info := taskGroupRepository.ToEntity(taskGroupName)
	item := mapper.Single[response.TaskGroupResponse](info)
	lstClient := clientRepository.ToList()
	item.Clients = mapper.ToList[response.ClientResponse](getTaskGroupClientList(item.Name, lstClient))
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

	// 发到所有节点上
	_ = container.Resolve[core.IEvent]("TaskGroupUpdate").Publish(taskGroupDO)
}

// 任务组删除
// @post delete
func TaskGroupDelete(taskGroupName string, taskGroupRepository taskGroup.Repository) {
	// 判断任务组是否存在
	exception.ThrowWebExceptionBool(!taskGroupRepository.IsExists(taskGroupName), 403, "任务组不存在")
	taskGroupRepository.Delete(taskGroupName)
	domain.RemoveMonitorTaskGroup(taskGroupName)
}

// 设置任务组状态
// @post setEnable
func SetEnable(taskGroupName string, enable bool, taskGroupRepository taskGroup.Repository) {
	taskGroupDO := taskGroupRepository.ToEntity(taskGroupName)
	taskGroupDO.SetEnable(enable)
	taskGroupRepository.Save(taskGroupDO)
	// 发到所有节点上
	_ = container.Resolve[core.IEvent]("TaskGroupUpdate").Publish(taskGroupDO)
}

func getTaskGroupClientList(taskGroupName string, lstClient collections.List[client.DomainObject]) collections.List[client.DomainObject] {
	// 筛选包含任务组的客户端
	lstClient = lstClient.Where(func(item client.DomainObject) bool {
		return item.Job.Name == taskGroupName
	}).ToList()

	return lstClient.OrderByDescending(func(item client.DomainObject) any {
		return item.IsMaster
	}).ToList()
}
