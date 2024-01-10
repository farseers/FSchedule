// @area /api/
package api

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/exception"
)

// TaskReport 客户端回调
// @post /taskReport
func TaskReport(dto client.TaskReportVO, taskGroupRepository taskGroup.Repository, scheduleRepository schedule.Repository) {
	switch dto.Status {
	case enum.None, enum.Scheduling, enum.ScheduleFail:
		exception.ThrowWebExceptionf(403, "任务组 %s %d 回调的状态设置不正确：%s", dto.Name, dto.Id, dto.Status.String())
	case enum.Working, enum.Fail, enum.Success: // 正确的，继续执行
	}

	// 加锁
	scheduleRepository.ScheduleLock(dto.Name, dto.Id).GetLockRun(func() {
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
			// 仅更新任务
			taskEO.UpdateTask(dto.Status, dto.Data, dto.Progress, dto.RunSpeed)
			taskGroupRepository.SaveTask(taskEO)
			return
		}

		// 更新任务
		taskGroupDO.Report(dto.Status, dto.Data, dto.Progress, dto.RunSpeed, dto.NextTimespan, taskGroupRepository)
	})
}
