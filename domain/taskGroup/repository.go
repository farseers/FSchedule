package taskGroup

type Repository interface {
	// ToEntity 获取任务组信息
	ToEntity(name string) DomainObject
	// Save 保存任务组信息
	Save(do DomainObject)
}
