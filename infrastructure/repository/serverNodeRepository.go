package repository

import (
	"FSchedule/domain/serverNode"
	"FSchedule/infrastructure/repository/context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"strconv"
)

const serverCacheKey = "FSchedule_ServerNode"

type serverNodeRepository struct {
}

func (receiver *serverNodeRepository) Save(do *serverNode.DomainObject) {
	if do.Id == 0 {
		return
	}
	do.ActivateAt = dateTime.Now()
	_ = context.RedisContext("更新服务端节点").HashSetEntity(serverCacheKey, strconv.FormatInt(do.Id, 10), &do)
}

func (receiver *serverNodeRepository) ToList() collections.List[serverNode.DomainObject] {
	var servers []serverNode.DomainObject
	_ = context.RedisContext("获取服务端节点列表").HashToArray(serverCacheKey, &servers)
	return collections.NewList(servers...)
}

func (receiver *serverNodeRepository) Remove(id int64) {
	_, _ = context.RedisContext("移除服务端节点列表").HashDel(serverCacheKey, strconv.FormatInt(id, 10))
}

func (receiver *serverNodeRepository) GetCount() int64 {
	count := context.RedisContext("获取服务端节点数量").HashCount(serverCacheKey)
	return int64(count)
}

func (receiver *serverNodeRepository) ToEntity(serverId int64) serverNode.DomainObject {
	var do serverNode.DomainObject
	_, _ = context.RedisContext("获取服务端节点").HashToEntity(serverCacheKey, strconv.FormatInt(serverId, 10), &do)
	return do
}
