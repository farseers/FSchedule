package client

import "github.com/farseer-go/collections"

type Repository interface {
	// Save 保存客户端信息
	Save(do DomainObject)
	// ToList 获取客户端列表
	ToList() collections.List[DomainObject]
	// RemoveClient 移除客户端
	RemoveClient(clientId string)
	// GetCount 客户端数量
	GetCount() int64
	// ToEntity 获取客户端
	ToEntity(clientId string) DomainObject
	// 同步客户端
	Sync(lst collections.List[DomainObject])
}
