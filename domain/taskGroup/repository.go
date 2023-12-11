package taskGroup

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
)

type Repository interface {
	taskRepository
	// ToEntity 获取任务组信息
	ToEntity(id int64) DomainObject
	// ToList 获取所有任务组中的任务
	ToList() collections.List[DomainObject]
	// ToListByName 获取所有任务组中的任务
	ToListByName(name string) collections.List[DomainObject]
	// Save 保存任务组信息
	Save(do DomainObject)
	// SaveAndTask 保存任务组、任务信息
	SaveAndTask(do DomainObject)
	// SaveTask 保存任务信息
	SaveTask(taskEO TaskEO)
	// GetTask 获取任务信息
	GetTask(taskGroupId int64, taskId int64) TaskEO
	// Sync 同步任务组数据
	Sync()

	// *******************仪表盘使用*********************
	ToListForPage(name string, enable int, taskStatus enum.TaskStatus, clientId int64, pageSize int, pageIndex int) collections.PageList[DomainObject]
	IsExists(taskGroupId int64) bool                                                       // 任务组是否存在
	UpdateByEdit(do DomainObject)                                                          // 修改
	Delete(taskGroupId int64)                                                              // 删除
	GetTaskGroupCount() int64                                                              // 任务组数量
	GetUnRunCount() int                                                                    // 超时未运行的任务组数量
	ToSchedulerWorkingList(pageSize int, pageIndex int) collections.PageList[DomainObject] // 调度中的任务组
	GetUnRunList(pageSize int, pageIndex int) collections.PageList[DomainObject]           // 超时未运行的任务组
}

type taskRepository interface {
	// ToTaskSpeedList 当前任务组下所有任务的执行速度
	ToTaskSpeedList(taskGroupId int64) []int64
	// ClearFinish 清除成功的任务记录（1天前）
	TaskClearFinish(taskGroupId int64, taskId int)
	ToTaskFinishList(taskGroupId int64, top int) collections.List[TaskEO]
	// *******************仪表盘使用*********************
	// ToFinishList 获取指定任务组执行成功的任务列表
	ToTaskListByGroupId(taskGroupId int64, taskStatus enum.TaskStatus, pageSize int, pageIndex int) collections.PageList[TaskEO]
	TodayFailCount() int64 // 今天失败数量
}
