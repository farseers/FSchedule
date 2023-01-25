package repository

import (
	"FSchedule/domain/client"
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
