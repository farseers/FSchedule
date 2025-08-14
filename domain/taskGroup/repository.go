package taskGroup

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"

	"github.com/farseer-go/collections"
)

type Repository interface {
	TaskRepository
	// ToEntity 获取任务组信息
	ToEntity(taskGroupName string) DomainObject
	// ToList 获取所有任务组中的任务
	ToList() collections.List[DomainObject]
	// Save 保存任务组信息
	Save(do DomainObject)
	// SaveAndTask 保存任务组、任务信息
	SaveAndTask(do DomainObject)
	// Sync 同步任务组数据
	Sync()

	// *******************仪表盘使用*********************
	ToListForFops(taskGroupName string, enable int, taskStatus executeStatus.Enum, taskId int64, clientId string, pageSize int, pageIndex int) collections.List[DomainObject]
	IsExists(taskGroupName string) bool // 任务组是否存在
	Delete(taskGroupName string)        // 删除
	GetTaskGroupCount() int64           // 任务组数量
	GetUnRunCount() int                 // 超时未运行的任务组数量
}

type TaskRepository interface {
	ToTaskSpeedList() collections.List[TaskEO]                           // ToTaskSpeedList 当前任务组下所有任务的执行速度
	TaskClearFinish(taskGroupName string, taskId int64)                  // TaskClearFinish 清除成功的任务记录（1天前）
	GetLastFinishTaskId(reservedTaskCount int) (map[string]int64, error) // 获取已完成的任务TaskId
	// *******************仪表盘使用*********************
	// ToTaskListByGroupId 获取指定任务组执行成功的任务列表
	ToHistoryTaskList(clientName, taskGroupName string, scheduleStatus scheduleStatus.Enum, executeStatus executeStatus.Enum, taskId string, pageSize int, pageIndex int) collections.PageList[TaskEO]
	TodayFailCount() int64                      // 今天失败数量
	GetStatCount() collections.List[StatTaskEO] // 统计任务成功失败数量
	// SaveTask 保存任务信息
	SaveTask(taskEO TaskEO)
	// GetTask 获取任务信息
	GetTask(taskGroupName string, taskId int64) TaskEO
}
