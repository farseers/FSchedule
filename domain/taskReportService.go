package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/taskGroup"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
)

func TaskReportService(clientId string, dto client.TaskReportVO, taskGroupRepository taskGroup.Repository) {
	// 1. 获取 Master 对象（判断当前节点是不是指挥官）
	taskGroupMonitor := GetTaskGroupMonitorMaster(dto.Name)

	var taskGroupDO *taskGroup.DomainObject
	if taskGroupMonitor != nil {
		taskGroupDO = taskGroupMonitor.DomainObject
	} else {
		// 如果本地不是 Master，从数据库取最新的实体（保证数据更新的基准是对的）
		do := taskGroupRepository.ToEntity(dto.Name)
		if do.IsNil() {
			return
		}
		taskGroupDO = &do
	}

	// 2. 处理 ID 不匹配的情况（旧任务补报）
	if taskGroupDO.Task.Id != dto.Id {
		taskEO := taskGroupRepository.GetTask(dto.Name, dto.Id)
		if !taskEO.IsNull() {
			taskEO.UpdateTask(dto.Status, dto.Data, dto.Progress, dto.FailRemark)
			taskGroupRepository.SaveTask(taskEO)
		}
		return
	}

	// 3. 正常状态报备
	taskGroupDO.Report(dto.Status, dto.Data, dto.Progress, dto.NextTimespan, dto.FailRemark, taskGroupRepository)
	taskGroupRepository.Save(*taskGroupDO)

	// 4. 通知决策
	if taskGroupMonitor != nil {
		taskGroupMonitor.Notify()
		return
	}

	// 说明master不在当前节点,则广播出去
	_ = container.Resolve[core.IEvent]("TaskGroupUpdate").Publish(*taskGroupDO)
}
