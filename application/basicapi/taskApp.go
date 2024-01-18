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

// 今天失败数量
// @get todayFailCount
func TodayFailCount(taskGroupRepository taskGroup.Repository) int64 {
	return taskGroupRepository.TodayFailCount()
}

// 统计异常数量列表
// @get statList
func StatList(taskGroupRepository taskGroup.Repository) collections.List[taskGroup.StatTaskEO] {
	return taskGroupRepository.GetStatCount()
}
