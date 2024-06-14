// @area /basicapi/stat/
package basicapi

import (
	"FSchedule/application/basicapi/response"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
)

// 统计异常数量列表
// @get statList
func StatList(taskGroupRepository taskGroup.Repository) collections.List[taskGroup.StatTaskEO] {
	return taskGroupRepository.GetStatCount()
}

// 统计数量
// @get info
func Info(taskGroupRepository taskGroup.Repository) response.InfoResponse {
	return response.InfoResponse{
		TaskGroupCount:      taskGroupRepository.GetTaskGroupCount(),
		TaskGroupUnRunCount: taskGroupRepository.GetUnRunCount(),
		TodayFailCount:      taskGroupRepository.TodayFailCount(),
	}
}
