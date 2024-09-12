// @area /basicapi/client/
package basicapi

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/collections"
)

// 客户端列表
// @get list
func ClientList(clientRepository client.Repository) collections.List[client.DomainObject] {
	lst := clientRepository.ToList().OrderBy(func(item client.DomainObject) any {
		return item.Job.Name + item.Id
	}).ToList()
	return lst
}
