// @area /api/
package api

import (
	"FSchedule/domain/client"
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
)

// TaskReport 客户端回调
// @post taskReport
func TaskReport(dto client.TaskReportVO, taskGroupRepository taskGroup.Repository, scheduleRepository schedule.Repository) {
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
			taskEO.UpdateTask(dto.Status, dto.Data, dto.Progress, dto.FailRemark)
			taskGroupRepository.SaveTask(taskEO)
			return
		}

		// 更新任务
		if taskGroupDO.Name != taskGroupDO.Task.Name {
			_ = flog.Errorf("任务组：%s 注意，客户端回调，发现task.Name不一致，TaskId=%d，taskName=%s, task=%+v", taskGroupDO.Name, taskGroupDO.Task.Id, taskGroupDO.Task.Name, taskGroupDO.Task)
		}
		taskGroupDO.Report(dto.Status, dto.Data, dto.Progress, dto.NextTimespan, dto.FailRemark, taskGroupRepository)
	})
}

// TaskReport Kill任务
// @post killTask
func KillTask(taskGroupName string, taskGroupRepository taskGroup.Repository, clientRepository client.Repository, clientCheck client.IClientCheck) {
	taskGroupDO := taskGroupRepository.ToEntity(taskGroupName)
	if taskGroupDO.IsNil() {
		exception.ThrowWebExceptionf(403, "任务组 %s 不存在", taskGroupName)
	}

	if taskGroupDO.Task.IsFinish() {
		exception.ThrowWebExceptionf(403, "任务组 %s %d 状态为已完成，无法停止。", taskGroupDO.Name, taskGroupDO.Task.Id)
	}

	// 通知客户端，停止任务
	if taskGroupDO.Task.Client.Id > 0 {
		clientDO := clientRepository.ToEntity(taskGroupDO.Task.Client.Id)
		if err := clientCheck.Kill(clientDO, taskGroupDO.Task.Id); err != nil {
			flog.Warningf("任务组 %s %d 通知客户端%s（%d）：%s:%d 停止任务失败：%s", taskGroupDO.Name, taskGroupDO.Task.Id, clientDO.Name, clientDO.Id, clientDO.Ip, clientDO.Port, err.Error())
		}
	}

	// 更新任务状态
	if taskGroupDO.Name != taskGroupDO.Task.Name {
		_ = flog.Errorf("任务组：%s 注意，Kill任务，发现task.Name不一致，TaskId=%d，taskName=%s, task=%+v", taskGroupDO.Name, taskGroupDO.Task.Id, taskGroupDO.Task.Name, taskGroupDO.Task)
	}
	taskGroupDO.Report(executeStatus.Fail, taskGroupDO.Data, taskGroupDO.Task.Progress, 0, "FOPS主动停止", taskGroupRepository)
}
