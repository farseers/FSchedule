package client

import (
	"sync"
)

var clientList = sync.Map{}

func GetClient(clientId string) *DomainObject {
	if clientDO, exists := clientList.Load(clientId); exists {
		return clientDO.(*DomainObject)
	}
	return nil
}
