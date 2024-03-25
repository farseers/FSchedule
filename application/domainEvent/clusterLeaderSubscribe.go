package domainEvent

import (
	"FSchedule/application/job"
	"FSchedule/domain"
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
	if leaderId == core.AppId {
		// 更新集群leader信息
		serverNodeRepository := container.Resolve[serverNode.Repository]()
		lst := serverNodeRepository.ToList()
		for i := 0; i < lst.Count(); i++ {
			serverNodeDO := lst.Index(i)
			serverNodeDO.SetLeader(leaderId)
			serverNodeRepository.Save(&serverNodeDO)
		}

		// 同步任务组、任务数据
		if configure.GetInt("FSchedule.DataSyncTime") > 0 {
			tasks.Run("taskGroupSync", 10*time.Minute, func(context *tasks.TaskContext) {
				container.Resolve[taskGroup.Repository]().Sync()
			}, fs.Context)
		}

		// 标记当前节点为Leader
		domain.CheckOnline()
		serverNode.IsLeaderNode = true

		// 移除30秒不活跃的
		tasks.Run("ServerNodeTimeoutJob", 30*time.Second, job.ServerNodeTimeoutJob, fs.Context)

		// 计算任务组的平均耗时
		tasks.Run("SyncAvgSpeedJob", 30*time.Minute, job.SyncAvgSpeedJob, fs.Context)

		// 自动清除历史任务记录
		if configure.GetInt("FSchedule.ReservedTaskCount") > 0 {
			tasks.Run("ClearHisTaskJob", 1*time.Hour, job.ClearHisTaskJob, fs.Context)
		}
	}
}
