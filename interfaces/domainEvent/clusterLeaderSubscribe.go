package domainEvent

import (
	"FSchedule/domain/serverNode"
	"FSchedule/domain/taskGroup"
	"FSchedule/interfaces/job"
	"context"
	"time"

	"github.com/farseer-go/fs/color"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/redis"
	"github.com/farseer-go/tasks"
)

// ClusterLeaderSubscribe 选举事件
func ClusterLeaderSubscribe(message any, _ core.EventArgs) {
	leaderId := parse.Convert(message, int64(0))

	flog.Infof("选举%s为Master节点", color.Red(leaderId))

	// 当前节点是leader
	if leaderId == core.AppId {
		// 更新集群leader信息
		serverNodeRepository := container.Resolve[serverNode.Repository]()
		lst := serverNodeRepository.ToList()

		// redis事务更新
		_ = container.Resolve[redis.IClient]("default").Transaction(func() {
			lst.Foreach(func(item *serverNode.DomainObject) {
				item.SetLeader(leaderId)
				serverNodeRepository.Save(item)
			})
		})

		// 同步任务组、任务数据
		if sec := configure.GetInt("FSchedule.DataSyncTime"); sec > 0 {
			tasks.Run("taskGroupSync", time.Duration(sec)*time.Second, func(context *tasks.TaskContext) {
				// 使用事务同步任务到数据库
				container.Resolve[taskGroup.Repository]().Sync()
			}, context.Background())
		}

		serverNode.IsLeaderNode = true

		// 移除30秒不活跃的
		tasks.Run("ServerNodeTimeoutJob", 30*time.Second, job.ServerNodeTimeoutJob, context.Background())

		// 计算任务组的平均耗时（已废弃：改为每个TaskGroupMonitor独立异步计算，避免并发覆盖问题）
		// tasks.Run("SyncAvgSpeedJob", 30*time.Minute, job.SyncAvgSpeedJob, context.Background())

		// 自动清除历史任务记录
		if configure.GetInt("FSchedule.ReservedTaskCount") > 0 {
			tasks.RunNow("ClearHisTaskJob", 1*time.Hour, job.ClearHisTaskJob, context.Background())
		}

		// 每5秒检查客户端是否永久离线
		tasks.Run("SyncClientJob", 5*time.Second, job.RemoveClientJob, context.Background())
	}
}
