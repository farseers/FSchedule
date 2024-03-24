// @area /basicapi/task/
package basicapi

import (
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
)

// 任务列表
// @get list
func TaskList(clientName, taskGroupName string, taskStatus enum.TaskStatus, taskId int64, pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[taskGroup.TaskEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return taskGroupRepository.ToTaskListByGroupId(clientName, taskGroupName, taskStatus, taskId, pageSize, pageIndex)
}

// 按计划执行时间排序
// @get planList
func TaskPlanList(top int, taskGroupRepository taskGroup.Repository) collections.List[taskGroup.TaskEO] {
	lst := taskGroupRepository.ToList()
	// 先取任务
	var lstTask collections.List[taskGroup.TaskEO]
	lst.Select(&lstTask, func(item taskGroup.DomainObject) any {
		return item.Task
	})

	// 按时间排序
	return lstTask.OrderBy(func(item taskGroup.TaskEO) any {
		return item.StartAt.UnixMilli()
	}).Take(top).ToList()
}
