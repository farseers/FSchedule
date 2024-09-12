package repository

import (
	"FSchedule/domain/client"
	"FSchedule/infrastructure/repository/context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/redis"
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

func (receiver *clientRepository) GetCount() int64 {
	count := context.RedisContext("获取客户端数量").HashCount(clientCacheKey)
	return int64(count)
}

func (receiver *clientRepository) Sync(lst collections.List[client.DomainObject]) {
	_ = container.Resolve[redis.IClient]("default").Transaction(func() {
		_, _ = context.RedisContext("清除客户端").Del(clientCacheKey)
		lst.Foreach(func(clientDO *client.DomainObject) {
			_ = context.RedisContext("同步客户端").HashSetEntity(clientCacheKey, clientDO.Id, clientDO)
		})
	})
}
