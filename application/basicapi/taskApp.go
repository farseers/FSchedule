// @area /basicapi/task/
package basicapi

import (
	"FSchedule/application/basicapi/response"
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/mapper"
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
func TaskPlanList(top int, taskGroupRepository taskGroup.Repository) collections.List[response.TaskPlanResponse] {
	if top == 0 {
		top = 20
	}

	lst := taskGroupRepository.ToList()
	// 先取任务
	var lstTask collections.List[taskGroup.TaskEO]
	lst.Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable
	}).Select(&lstTask, func(item taskGroup.DomainObject) any {
		return item.Task
	})

	// 按时间排序
	lstTask = lstTask.OrderBy(func(item taskGroup.TaskEO) any {
		return item.StartAt.UnixMilli()
	}).Take(top).ToList()

	return mapper.ToList[response.TaskPlanResponse](lstTask, func(r *response.TaskPlanResponse, source any) {
		startAt := source.(taskGroup.TaskEO).StartAt
		isAfter := dateTime.Now().After(startAt)

		switch r.Status {
		case enum.None:
			if isAfter {
				r.StartAt = fmt.Sprintf("等待%s", dateTime.Now().Sub(startAt).String())
			} else {
				r.StartAt = fmt.Sprintf("超时%s", startAt.Sub(dateTime.Now()).String())
			}
		case enum.Scheduling, enum.Working:
			r.StartAt = fmt.Sprintf("已执行了%s", dateTime.Now().Sub(startAt).String())
		default:
		}
	})
}
