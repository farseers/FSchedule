package client

type IClientCheck interface {
	// Check 检查客户端存活
	Check(do *DomainObject) (ResourceVO, error)
	// Invoke 下发任务
	Invoke(do *DomainObject, task *TaskEO) (ResourceVO, error)
	// Status 查询任务状态
	Status(do *DomainObject, taskId int64) (TaskReportVO, error)
	// Kill 终止任务
	Kill(do DomainObject, taskId int64) bool
}
