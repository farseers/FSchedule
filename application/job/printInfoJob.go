package job

import (
	"FSchedule/domain"
	"FSchedule/domain/serverNode"
	"fmt"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
	"strings"
)

// PrintInfoJob 打印客户端、任务组信息
func PrintInfoJob(context *tasks.TaskContext) {
	lst := container.Resolve[serverNode.Repository]().ToList()
	var serverNodes []string
	lst.Select(&serverNodes, func(item serverNode.DomainObject) any {
		return fmt.Sprintf("%s %s:%d", flog.Green(item.Id), flog.Yellow(item.Ip), item.Port)
	})
	flog.Infof("%s个集群节点：%s", flog.Red(lst.Count()), strings.Join(serverNodes, ","))
	flog.Infof("%s个客户端（%s个正常状态），%s个任务组在监控（%s个在运行）", flog.Red(domain.ClientCount()), flog.Green(domain.ClientNormalCount()), flog.Red(domain.TaskGroupCount()), flog.Green(domain.TaskGroupEnableCount()))
}
