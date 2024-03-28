// @area /basicapi/task/
package basicapi

import (
	"FSchedule/application/basicapi/response"
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
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
func TaskList(clientName, taskGroupName string, scheduleStatus scheduleStatus.Enum, executeStatus executeStatus.Enum, taskId string, pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[taskGroup.TaskEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return taskGroupRepository.ToHistoryTaskList(clientName, taskGroupName, scheduleStatus, executeStatus, taskId, pageSize, pageIndex)
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
		return item.Name
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

		switch r.ExecuteStatus {
		case executeStatus.None:
			if r.ScheduleStatus != scheduleStatus.None {
				r.Plan = r.ScheduleStatus.String()
			}

			if isAfter {
				// 等待
				r.Plan += "等待 " + (time.Duration(startAt.Sub(dateTime.Now()).Seconds()) * time.Second).String()
			} else {
				// 超时
				r.Plan += "超时 " + (time.Duration(dateTime.Now().Sub(startAt).Seconds()) * time.Second).String()
			}
		case executeStatus.Working:
			r.Plan = fmt.Sprintf("已执行 %s", (time.Duration(dateTime.Now().Sub(schedulerAt).Seconds()) * time.Second).String())
		default:
		}

		r.Plan = strings.ReplaceAll(r.Plan, "m", "分")
		r.Plan = strings.ReplaceAll(r.Plan, "s", "秒")
		r.Plan = strings.ReplaceAll(r.Plan, "h", "时")
	})
}
