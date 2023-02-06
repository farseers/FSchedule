package domain

import "FSchedule/domain/taskGroup"

func TaskReportService(taskEO taskGroup.TaskEO, taskGroupDO taskGroup.DomainObject, taskGroupRepository taskGroup.Repository) {
	// 如果任务组的任务ID与任务ID不同时，则要单独保存任务数据
	if taskGroupDO.Task.Id != taskEO.Id {
		taskGroupRepository.SaveTask(taskEO)
	}

	// 更新任务组
	taskGroupDO.UpdateTask(taskEO)
	taskGroupRepository.Save(taskGroupDO)
}
