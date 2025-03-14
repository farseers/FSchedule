// @area /basicapi/client/
package basicapi

import (
	"FSchedule/domain/client"
	"strconv"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
)

// 客户端列表
// @get list
func ClientList(clientRepository client.Repository) collections.List[client.DomainObject] {
	lst := clientRepository.ToList().OrderBy(func(item client.DomainObject) any {
		return item.Job.Name + strconv.Itoa(parse.ToInt(item.IsMaster)) + item.Id
	}).ToList()
	return lst
}
