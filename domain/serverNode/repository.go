package serverNode

import (
	"github.com/farseer-go/collections"
)

type Repository interface {
	// Save 保存节点信息
	Save(do *DomainObject)
	// ToList 查询所有信息
	ToList() collections.List[DomainObject]
	// Remove 移除节点
	Remove(id int64)
	// GetCount 获取节点数量
	GetCount() int64
	// ToEntity 查询节点信息
	ToEntity(serverId int64) DomainObject
}
