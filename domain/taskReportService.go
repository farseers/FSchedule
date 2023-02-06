package domain

import "FSchedule/domain/taskGroup"

func TaskReportService(taskEO taskGroup.TaskEO, taskGroupDO taskGroup.DomainObject, taskGroupRepository taskGroup.Repository) {
	// 更新任务组
	taskGroupDO.UpdateTask(taskEO)
	taskGroupRepository.SaveAndTask(taskGroupDO)
}
