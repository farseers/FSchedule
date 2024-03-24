// @area /basicapi/stat/
package basicapi

import (
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
)

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
