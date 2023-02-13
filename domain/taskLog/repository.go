package taskLog

type Repository interface {
	// Add 添加日志
	Add(taskLogDO DomainObject)
}
