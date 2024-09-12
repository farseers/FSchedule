package repository

import (
	"FSchedule/domain/client"
	"FSchedule/infrastructure/repository/context"
	"github.com/farseer-go/collections"
)

const clientCacheKey = "FSchedule_ClientList"

type clientRepository struct {
}

func (receiver *clientRepository) ToEntity(clientId string) client.DomainObject {
	var do client.DomainObject
	_, _ = context.RedisContext("获取客户端").HashToEntity(clientCacheKey, clientId, &do)
	return do
}

func (receiver *clientRepository) Save(do client.DomainObject) {
	if do.Id == "" {
		return
	}
	_ = context.RedisContext("保存客户端").HashSetEntity(clientCacheKey, do.Id, do)
}

func (receiver *clientRepository) ToList() collections.List[client.DomainObject] {
	var clients []client.DomainObject
	_ = context.RedisContext("获取客户端列表").HashToArray(clientCacheKey, &clients)
	lst := collections.NewList(clients...)
	return lst.OrderBy(func(item client.DomainObject) any {
		return int(item.Status)
	}).ToList()
}

func (receiver *clientRepository) RemoveClient(clientId string) {
	_, _ = context.RedisContext("移除客户端").HashDel(clientCacheKey, clientId)
}
