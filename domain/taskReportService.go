package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/flog"
)

// TaskReport 客户端回调
func TaskReportService(clientId string, dto client.TaskReportVO, taskGroupRepository taskGroup.Repository) {
	var taskGroupDO *taskGroup.DomainObject
	taskGroupMonitor := GetTaskGroupMonitor(clientId)
	if taskGroupMonitor != nil {
		taskGroupDO = taskGroupMonitor.DomainObject
	} else {
		do := taskGroupRepository.ToEntity(dto.Name)
		if do.IsNil() {
			flog.Warningf("任务组[%s] 不存在", dto.Name)
			return
		}
		taskGroupDO = &do
	}

	// 任务ID不相同，说明任务上报，晚于当前任务组的最新任务，这时只要保存任务就可以了
	if taskGroupDO.Task.Id != dto.Id {
		taskEO := taskGroupRepository.GetTask(dto.Name, dto.Id)
		if taskEO.IsNull() {
			flog.Warningf("任务id={%d} 不存在", dto.Id)
			return
		}
		// 仅更新任务
		taskEO.UpdateTask(dto.Status, dto.Data, dto.Progress, dto.FailRemark)
		taskGroupRepository.SaveTask(taskEO)
		return
	}

	// 更新任务
	if taskGroupDO.Name != taskGroupDO.Task.Name {
		_ = flog.Errorf("任务组：%s 注意，客户端回调，发现task.Name不一致，TaskId=%d，taskName=%s, task=%+v", taskGroupDO.Name, taskGroupDO.Task.Id, taskGroupDO.Task.Name, taskGroupDO.Task)
	}
	taskGroupDO.Report(dto.Status, dto.Data, dto.Progress, dto.NextTimespan, dto.FailRemark, taskGroupRepository)
	// 为了让FOPS可以立即查询到工作状态，这里需要立即保存
	if dto.Status == executeStatus.Working {
		taskGroupRepository.Save(*taskGroupDO)
	}
	if taskGroupMonitor != nil {
		taskGroupMonitor.updated <- struct{}{}
	}
}
