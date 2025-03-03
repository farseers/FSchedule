package job

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/serverNode"
	"fmt"
	"strings"

	"github.com/farseer-go/fs/color"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/tasks"
)

// PrintInfoJob 打印客户端、任务组信息
func PrintInfoJob(context *tasks.TaskContext) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	clientCount := client.GetClientCount()
	if clientCount > 0 {
		flog.Printf("%s个客户端，%s个任务组在运行（%s个连接）\n", color.Red(clientCount), color.Red(domain.TaskGroupCount()), color.Green(domain.TaskGroupEnableCount()))
	}
	lst := container.Resolve[serverNode.Repository]().ToList()

	// 主节点
	leaderNode := lst.Where(func(item serverNode.DomainObject) bool {
		return item.IsLeader
	}).OrderByDescending(func(item serverNode.DomainObject) any {
		return item.ActivateAt
	}).First()
	if leaderNode.Id > 0 {
		flog.Printf("%s个Master节点：%s %s:%d\n", color.Red(1), color.Green(leaderNode.Id), color.Yellow(leaderNode.Ip), leaderNode.Port)
	}

	// 从节点
	var serverNodes []string
	lst.Where(func(item serverNode.DomainObject) bool {
		return !item.IsLeader
	}).OrderByDescending(func(item serverNode.DomainObject) any {
		return item.Id
	}).Select(&serverNodes, func(item serverNode.DomainObject) any {
		return fmt.Sprintf("%s %s:%d", color.Blue(item.Id), item.Ip, item.Port)
	})
	if len(serverNodes) > 0 {
		flog.Printf("%s个Slave节点：%s\n", color.Red(len(serverNodes)), strings.Join(serverNodes, ","))
	}
	flog.Println("---------------------------------------")
}
