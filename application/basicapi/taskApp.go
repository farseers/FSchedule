// @area /basicapi/task/
package basicapi

import (
	"FSchedule/domain/enum"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
)

// 任务列表
// @get list
func TaskList(taskGroupName string, taskStatus enum.TaskStatus, pageSize int, pageIndex int, taskGroupRepository taskGroup.Repository) collections.PageList[taskGroup.TaskEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return taskGroupRepository.ToTaskListByGroupId(taskGroupName, taskStatus, pageSize, pageIndex)
}

// 今天失败数量
// @get todayFailCount
func TodayFailCount(taskGroupRepository taskGroup.Repository) int64 {
	return taskGroupRepository.TodayFailCount()
}
