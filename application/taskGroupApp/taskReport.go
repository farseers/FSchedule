package taskGroupApp

import (
	"FSchedule/domain/client"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/exception"
)

// TaskReport 客户端回调
func TaskReport(dto client.TaskReportVO, taskGroupRepository taskGroup.Repository, scheduleRepository schedule.Repository) {
	// 加锁
	scheduleRepository.NewLock(dto.Name).GetLockRun(func() {
		taskGroupDO := taskGroupRepository.ToEntity(dto.Name)
		if taskGroupDO.IsNil() {
			exception.ThrowWebExceptionf(403, "任务组[%s] 不存在", dto.Name)
		}
		// 任务ID不相同，说明任务上报，晚于当前任务组的最新任务，这时只要保存任务就可以了
		if taskGroupDO.Task.Id != dto.Id {
			taskEO := taskGroupRepository.GetTask(dto.Name, dto.Id)
			if taskEO.IsNull() {
				exception.ThrowWebExceptionf(403, "任务id={%d} 不存在", dto.Id)
			}
			// 更新任务
			taskEO.UpdateTask(dto.Status, dto.Data, dto.Progress, dto.RunSpeed)
			taskGroupRepository.SaveTask(taskEO)
			return
		}

		taskGroupDO.Report(dto.Status, dto.Data, dto.Progress, dto.RunSpeed, dto.NextTimespan, taskGroupRepository)
	})
}
