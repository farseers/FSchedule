package repository

import (
	"FSchedule/domain/client"
	"FSchedule/infrastructure/repository/context"
	"github.com/farseer-go/collections"
	"strconv"
)

const clientCacheKey = "FSchedule_ClientList"

type clientRepository struct {
}

func (receiver *clientRepository) Save(do client.DomainObject) {
	if do.Id == 0 {
		return
	}
	_ = context.RedisContext("保存客户端").HashSetEntity(clientCacheKey, strconv.FormatInt(do.Id, 10), do)
}

func (receiver *clientRepository) ToList() collections.List[client.DomainObject] {
	var clients []client.DomainObject
	_ = context.RedisContext("获取客户端列表").HashToArray(clientCacheKey, &clients)
	lst := collections.NewList(clients...)
	return lst.OrderBy(func(item client.DomainObject) any {
		return int(item.Status)
	}).ToList()
}

func (receiver *clientRepository) RemoveClient(id int64) {
	_, _ = context.RedisContext("移除客户端").HashDel(clientCacheKey, strconv.FormatInt(id, 10))
}

func (receiver *clientRepository) GetCount() int64 {
	count := context.RedisContext("获取客户端数量").HashCount(clientCacheKey)
	return int64(count)
}

func (receiver *clientRepository) ToEntity(clientId int64) client.DomainObject {
	var do client.DomainObject
	_, _ = context.RedisContext("获取客户端").HashToEntity(clientCacheKey, strconv.FormatInt(clientId, 10), &do)
	return do
}
