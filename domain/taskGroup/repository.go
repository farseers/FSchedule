package taskGroup

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
)

type Repository interface {
	TaskRepository
	// ToEntity 获取任务组信息
	ToEntity(taskGroupName string) DomainObject
	// ToList 获取所有任务组中的任务
	ToList() collections.List[DomainObject]
	// ToListByName 获取所有任务组中的任务
	ToListByName(taskGroupName string) collections.List[DomainObject]
	// Save 保存任务组信息
	Save(do DomainObject)
	// SaveAndTask 保存任务组、任务信息
	SaveAndTask(do DomainObject)
	// GetTask 获取任务信息
	GetTask(taskGroupName string, taskId int64) TaskEO
	// Sync 同步任务组数据
	Sync()

	// *******************仪表盘使用*********************
	ToListForPage(clientName, taskGroupName string, enable int, taskStatus enum.TaskStatus, taskId, clientId int64, pageSize int, pageIndex int) collections.PageList[DomainObject]
	IsExists(taskGroupName string) bool                                                    // 任务组是否存在
	Delete(taskGroupName string)                                                           // 删除
	GetTaskGroupCount() int64                                                              // 任务组数量
	GetUnRunCount() int                                                                    // 超时未运行的任务组数量
	ToSchedulerWorkingList(pageSize int, pageIndex int) collections.PageList[DomainObject] // 调度中的任务组
	GetUnRunList(pageSize int, pageIndex int) collections.PageList[DomainObject]           // 超时未运行的任务组
}

type TaskRepository interface {
	// ToTaskSpeedList 当前任务组下所有任务的执行速度
	ToTaskSpeedList() collections.List[TaskEO]
	// TaskClearFinish 清除成功的任务记录（1天前）
	TaskClearFinish(taskGroupName string, taskId int64)
	ToTaskFinishList(taskGroupName string, top int) collections.List[TaskEO]
	// *******************仪表盘使用*********************
	// ToTaskListByGroupId 获取指定任务组执行成功的任务列表
	ToTaskListByGroupId(clientName, taskGroupName string, taskStatus enum.TaskStatus, taskId int64, pageSize int, pageIndex int) collections.PageList[TaskEO]
	TodayFailCount() int64 // 今天失败数量
	GetStatCount() collections.List[StatTaskEO]
	// SaveTask 保存任务信息
	SaveTask(taskEO TaskEO)
}
