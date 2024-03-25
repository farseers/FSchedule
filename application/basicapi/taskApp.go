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
	"strings"
	"time"
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
	// 状态为可用、非完成状态、并按开始时间排序
	lst = lst.Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && !item.Task.IsFinish()
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.Task.StartAt.UnixNano()
	}).Take(top).ToList()

	// 先取任务
	var lstTask collections.List[taskGroup.TaskEO]
	lst.Select(&lstTask, func(item taskGroup.DomainObject) any {
		return item.Task
	})

	return mapper.ToList[response.TaskPlanResponse](lstTask, func(r *response.TaskPlanResponse, source any) {
		startAt := source.(taskGroup.TaskEO).StartAt
		schedulerAt := source.(taskGroup.TaskEO).SchedulerAt
		isAfter := startAt.After(dateTime.Now())

		switch r.Status {
		case enum.None:
			if isAfter {
				r.StartAt = fmt.Sprintf("等待 %s", (time.Duration(startAt.Sub(dateTime.Now()).Seconds()) * time.Second).String())
			} else {
				r.StartAt = fmt.Sprintf("超时 %s", (time.Duration(dateTime.Now().Sub(startAt).Seconds()) * time.Second).String())
			}
		case enum.ScheduleFail:
			r.StartAt = fmt.Sprintf("调度失败，超时 %s", (time.Duration(dateTime.Now().Sub(startAt).Seconds()) * time.Second).String())
		case enum.Scheduling, enum.Working:
			r.StartAt = fmt.Sprintf("已执行 %s", (time.Duration(dateTime.Now().Sub(schedulerAt).Seconds()) * time.Second).String())
		default:
		}
		r.StartAt = strings.ReplaceAll(r.StartAt, "m", "分")
		r.StartAt = strings.ReplaceAll(r.StartAt, "s", "秒")
		r.StartAt = strings.ReplaceAll(r.StartAt, "h", "时")
	})
}
