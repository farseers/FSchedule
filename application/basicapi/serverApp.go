// @area /basicapi/server/
package basicapi

import (
	"FSchedule/domain/serverNode"
	"github.com/farseer-go/collections"
)

// 服务端节点列表
// @get list
func ServerList(serverNodeRepository serverNode.Repository) collections.List[serverNode.DomainObject] {
	return serverNodeRepository.ToList()
}
