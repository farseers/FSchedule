package repository

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/redis"
	"strconv"
)

const clientCacheKey = "FSS_ClientList"

type clientRepository struct {
	*redis.Client
}

func (receiver *clientRepository) Save(do client.DomainObject) {
	_ = receiver.Hash.SetEntity(clientCacheKey, strconv.FormatInt(do.Id, 10), &do)
}

func (receiver *clientRepository) ToList() collections.List[client.DomainObject] {
	var clients []client.DomainObject
	_ = receiver.Hash.ToArray(clientCacheKey, &clients)
	return collections.NewList(clients...)
}

func (receiver *clientRepository) RemoveClient(id int64) {
	_, _ = receiver.Hash.Del(clientCacheKey, strconv.FormatInt(id, 10))
}

func (receiver *clientRepository) GetCount() int64 {
	count := receiver.Hash.Count(clientCacheKey)
	return int64(count)
}

func (receiver *clientRepository) ToEntity(clientId int64) *client.DomainObject {
	var do *client.DomainObject
	_, _ = receiver.Hash.ToEntity(clientCacheKey, strconv.FormatInt(clientId, 10), do)
	return do
}
