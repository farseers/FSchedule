package repository

import (
	"FSchedule/domain/serverNode"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/redis"
	"strconv"
	"time"
)

const serverCacheKey = "FSchedule_ServerNode"

type serverNodeRepository struct {
	redis.IClient `inject:"default"`
}

func (receiver *serverNodeRepository) Save(do *serverNode.DomainObject) {
	if do.Id == 0 {
		return
	}
	do.ActivateAt = time.Now()
	_ = receiver.HashSetEntity(serverCacheKey, strconv.FormatInt(do.Id, 10), &do)
}

func (receiver *serverNodeRepository) ToList() collections.List[serverNode.DomainObject] {
	var servers []serverNode.DomainObject
	_ = receiver.HashToArray(serverCacheKey, &servers)

	for _, node := range servers {
		if time.Since(node.ActivateAt).Seconds() >= 30 {
			receiver.Remove(node.Id)
		}
	}
	return collections.NewList(servers...)
}

func (receiver *serverNodeRepository) Remove(id int64) {
	_, _ = receiver.HashDel(serverCacheKey, strconv.FormatInt(id, 10))
}

func (receiver *serverNodeRepository) GetCount() int64 {
	count := receiver.HashCount(serverCacheKey)
	return int64(count)
}

func (receiver *serverNodeRepository) ToEntity(serverId int64) serverNode.DomainObject {
	var do serverNode.DomainObject
	_, _ = receiver.HashToEntity(serverCacheKey, strconv.FormatInt(serverId, 10), &do)
	return do
}
