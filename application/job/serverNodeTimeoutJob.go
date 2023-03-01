package job

import (
	"FSchedule/domain/serverNode"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
	"time"
)

// ServerNodeTimeoutJob 移除30秒不活跃的
func ServerNodeTimeoutJob(context *tasks.TaskContext) {
	repository := container.Resolve[serverNode.Repository]()
	lst := repository.ToList()

	for i := 0; i < lst.Count(); i++ {
		serverNodeDO := lst.Index(i)
		if time.Since(serverNodeDO.ActivateAt).Seconds() >= 30 {
			repository.Remove(serverNodeDO.Id)
			flog.Infof("集群节点：%s %s:%d 不再活跃，移出集群", flog.Green(serverNodeDO.Id), flog.Yellow(serverNodeDO.Ip), serverNodeDO.Port)
		}
	}
}
