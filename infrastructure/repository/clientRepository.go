package repository

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/redis"
	"strconv"
)

const clientCacheKey = "FSS_ClientList"
const jobClientCacheKey = "FSS_JobClientList"

type clientRepository struct {
	*redis.Client
}

func (receiver *clientRepository) Save(do *client.DomainObject) {
	_ = receiver.Hash.SetEntity(clientCacheKey, strconv.FormatInt(do.Id, 10), &do)

	// 将客户端支持的任务列表保存到另外的KEY，方便通过任务名称来查找客户端列表
	// 这里可以用redis事务
	for _, job := range do.Jobs {
		key := fmt.Sprintf("%s:%s:%d", jobClientCacheKey, job.Name, job.Ver)
		_ = receiver.Hash.SetEntity(key, strconv.FormatInt(do.Id, 10), &do)
	}
	domain.MonitorClientPush(do)
}

func (receiver *clientRepository) ToList() collections.List[client.DomainObject] {
	var clients []client.DomainObject
	_ = receiver.Hash.ToArray(clientCacheKey, &clients)
	return collections.NewList(clients...)
}

// GetClients 获取支持taskGroupName的客户端列表
func (receiver *clientRepository) GetClients(taskGroupName string, version int) collections.List[client.DomainObject] {
	key := fmt.Sprintf("%s:%s:%d", jobClientCacheKey, taskGroupName, version)
	var clients []client.DomainObject
	_ = receiver.Hash.ToArray(key, &clients)
	return collections.NewList(clients...)
}

func (receiver *clientRepository) RemoveClient(id int64) {
	// 先移除客户端支持的任务
	clientDO := receiver.ToEntity(id)
	for _, job := range clientDO.Jobs {
		key := fmt.Sprintf("%s:%s:%d", jobClientCacheKey, job.Name, job.Ver)
		_, _ = receiver.Hash.Del(key, strconv.FormatInt(id, 10))
	}
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
