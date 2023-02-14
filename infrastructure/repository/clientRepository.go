package repository

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/redis"
	"strconv"
)

const clientCacheKey = "FSchedule_ClientList"

type clientRepository struct {
	redis.IClient        `inject:"default"`
	ClientUpdateEventBus core.IEvent `inject:"ClientUpdate"`
}

func (receiver *clientRepository) Save(do *client.DomainObject) {
	if do.Id == 0 {
		return
	}
	_ = receiver.HashSetEntity(clientCacheKey, strconv.FormatInt(do.Id, 10), &do)

	// 发到所有节点上
	_ = receiver.ClientUpdateEventBus.Publish(do)
}

func (receiver *clientRepository) ToList() collections.List[client.DomainObject] {
	var clients []client.DomainObject
	_ = receiver.HashToArray(clientCacheKey, &clients)
	return collections.NewList(clients...)
}

func (receiver *clientRepository) RemoveClient(id int64) {
	_, _ = receiver.HashDel(clientCacheKey, strconv.FormatInt(id, 10))
}

func (receiver *clientRepository) GetCount() int64 {
	count := receiver.HashCount(clientCacheKey)
	return int64(count)
}

func (receiver *clientRepository) ToEntity(clientId int64) client.DomainObject {
	var do client.DomainObject
	_, _ = receiver.HashToEntity(clientCacheKey, strconv.FormatInt(clientId, 10), &do)
	return do
}
