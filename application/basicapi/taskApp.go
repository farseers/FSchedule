// @area /basicapi/task/
package basicapi

import (
	"FSchedule/application/basicapi/response"
	"FSchedule/domain/client"
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
	"FSchedule/domain/taskGroup"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"strings"
	"time"
)

// 任务列表
// @get list
func TaskList(clientName, taskGroupName string, scheduleStatus scheduleStatus.Enum, executeStatus executeStatus.Enum, taskId string, pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[response.TaskResponse] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	lst := taskGroupRepository.ToHistoryTaskList(clientName, taskGroupName, scheduleStatus, executeStatus, taskId, pageSize, pageIndex)
	return mapper.ToPageList[response.TaskResponse](lst, func(r *response.TaskResponse, a any) {
		r.RunSpeed = (time.Duration(a.(taskGroup.TaskEO).RunSpeed) * time.Millisecond).String()
	})
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
				r.Plan = r.ScheduleStatus.String() + "，"
			}

			if isAfter {
				// 等待
				r.Plan += "等待 " + (time.Duration(startAt.Sub(dateTime.Now()).Seconds()) * time.Second).String()
				if r.Plan == "等待 0s" {
					r.Plan = "等待调度"
				}

			} else {
				// 超时
				r.Plan += "超时 " + (time.Duration(dateTime.Now().Sub(startAt).Seconds()) * time.Second).String()
			}
		case executeStatus.Working:
			r.Plan = fmt.Sprintf("已执行 %s", (time.Duration(dateTime.Now().Sub(schedulerAt).Seconds()) * time.Second).String())
		default:
		}

		r.Plan = strings.ReplaceAll(r.Plan, "等待 0s", "等待执行")
		r.Plan = strings.ReplaceAll(r.Plan, "m", "分")
		r.Plan = strings.ReplaceAll(r.Plan, "s", "秒")
		r.Plan = strings.ReplaceAll(r.Plan, "h", "时")
	})
}

// Kill任务
// @post killTask
func KillTask(taskGroupName string, taskGroupRepository taskGroup.Repository, clientRepository client.Repository, clientCheck client.IClientCheck) {
	taskGroupDO := taskGroupRepository.ToEntity(taskGroupName)
	if taskGroupDO.IsNil() {
		exception.ThrowWebExceptionf(403, "任务组 %s 不存在", taskGroupName)
	}

	if taskGroupDO.Task.IsFinish() {
		exception.ThrowWebExceptionf(403, "任务组 %s %d 状态为已完成，无法停止。", taskGroupDO.Name, taskGroupDO.Task.Id)
	}

	// 通知处理该任务组的服务端，需要调用客户端发起Kill请求
	taskGroupDO.Task.Kill = true
	taskGroupRepository.Save(taskGroupDO)

	// 发到所有节点上
	_ = container.Resolve[core.IEvent]("TaskGroupUpdate").Publish(taskGroupDO)
}
