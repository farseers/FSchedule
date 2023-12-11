// @area /basicapi/client/
package basicapi

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/collections"
)

// 客户端列表
// @get list
func ClientList(clientRepository client.Repository) collections.List[client.DomainObject] {
	return clientRepository.ToList()
}
