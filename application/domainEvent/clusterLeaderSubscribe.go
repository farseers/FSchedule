package domainEvent

import (
	"FSchedule/domain/serverNode"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/tasks"
	"time"
)

// ClusterLeaderSubscribe 选举事件
func ClusterLeaderSubscribe(message any, _ core.EventArgs) {
	leaderId := parse.Convert(message, int64(0))

	flog.Infof("选举%s为Master节点", flog.Red(leaderId))

	// 当前节点是leader
	if leaderId == fs.AppId {
		// 更新集群leader信息
		serverNodeRepository := container.Resolve[serverNode.Repository]()
		lst := serverNodeRepository.ToList()
		for i := 0; i < lst.Count(); i++ {
			serverNodeDO := lst.Index(i)
			serverNodeDO.SetLeader(leaderId)
			serverNodeRepository.Save(&serverNodeDO)
		}

		// 同步数据
		syncTime := configure.GetInt("FSchedule.DataSyncTime")
		if syncTime > 0 {
			// 每60s，同步一次任务组、任务
			tasks.RunNow("ServerNodeJob", 60*time.Second, func(context *tasks.TaskContext) {
				container.Resolve[taskGroup.Repository]().Sync()
			}, fs.Context)
		}
	}
}
