package repository

import (
	"FSchedule/domain/client"
	"FSchedule/infrastructure/repository/context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"strconv"
)

const clientCacheKey = "FSchedule_ClientList"

type clientRepository struct {
	ClientUpdateEventBus core.IEvent `inject:"ClientUpdate"`
}

func (receiver *clientRepository) Save(do *client.DomainObject) {
	if do.Id == 0 {
		return
	}
	_ = context.RedisContextIns.HashSetEntity(clientCacheKey, strconv.FormatInt(do.Id, 10), &do)

	// 发到所有节点上
	_ = receiver.ClientUpdateEventBus.Publish(do)
}

func (receiver *clientRepository) ToList() collections.List[client.DomainObject] {
	var clients []client.DomainObject
	_ = context.RedisContextIns.HashToArray(clientCacheKey, &clients)
	return collections.NewList(clients...)
}

func (receiver *clientRepository) RemoveClient(id int64) {
	_, _ = context.RedisContextIns.HashDel(clientCacheKey, strconv.FormatInt(id, 10))
}

func (receiver *clientRepository) GetCount() int64 {
	count := context.RedisContextIns.HashCount(clientCacheKey)
	return int64(count)
}

func (receiver *clientRepository) ToEntity(clientId int64) client.DomainObject {
	var do client.DomainObject
	_, _ = context.RedisContextIns.HashToEntity(clientCacheKey, strconv.FormatInt(clientId, 10), &do)
	return do
}
