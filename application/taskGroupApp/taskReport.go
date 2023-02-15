package taskGroupApp

import (
	"FSchedule/domain"
	"FSchedule/domain/enum"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
)

type TaskReportDTO struct {
	Id           int64                                  // 主键
	Name         string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Data         collections.Dictionary[string, string] // 数据
	NextTimespan int64                                  // 下次执行时间
	Progress     int                                    // 当前进度
	Status       enum.TaskStatus                        // 执行状态
	RunSpeed     int64                                  // 执行速度
}

// TaskReport 客户端回调
func TaskReport(dto TaskReportDTO, taskGroupRepository taskGroup.Repository, scheduleRepository schedule.Repository) {
	taskEO := taskGroupRepository.GetTask(dto.Name, dto.Id)
	if taskEO.IsNull() {
		exception.ThrowWebExceptionf(403, "任务id={%d} 不存在", dto.Id)
	}

	// 加锁
	scheduleRepository.NewLock(taskEO.Name).GetLockRun(func() {
		taskGroupDO := taskGroupRepository.ToEntity(taskEO.Name)
		if taskGroupDO.IsNil() {
			exception.ThrowWebExceptionf(403, "任务组[%s] 不存在", taskEO.Name)
		}

		// 重新计算下一个执行周期
		taskGroupDO.CalculateNextAtByUnix(dto.NextTimespan)
		// 更新任务
		taskEO.UpdateTask(dto.Status, dto.Data, dto.Progress, dto.RunSpeed)
		//flog.Debugf("任务组：%s 收到任务上报：%d", dto.Name, dto.Id)
		domain.TaskReportService(taskEO, taskGroupDO, taskGroupRepository)
	})
}
