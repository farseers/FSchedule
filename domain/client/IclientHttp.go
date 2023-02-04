package client

type IClientCheck interface {
	// Check 检查客户端是否存活
	Check(do *DomainObject) bool
	// Invoke 下发任务
	Invoke(do *DomainObject, task *TaskEO) bool
}
