package taskGroup

import "github.com/farseer-go/collections"

type Repository interface {
	// ToEntity 获取任务组信息
	ToEntity(name string) DomainObject
	// ToList 获取所有任务组中的任务
	ToList() collections.List[DomainObject]
	// Save 保存任务组信息
	Save(do DomainObject)
}
