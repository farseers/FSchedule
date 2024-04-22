package application

import (
	"FSchedule/application/job"
	"FSchedule/domain"
	"FSchedule/domain/schedule"
	"context"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/tasks"
	"time"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{domain.Module{}}
}

func (module Module) PostInitialize() {
	// 打印客户端、任务组信息
	fs.AddInitCallback("打印客户端、任务组信息", func() {
		tasks.Run("PrintInfoJob", 10*time.Second, job.PrintInfoJob, context.Background())
	})

	// 10秒更新一次服务端信息
	fs.AddInitCallback("每10秒更新节点活跃时间", func() {
		tasks.Run("ServerNodeJob", 10*time.Second, job.ServerActivateJob, context.Background())
	})

	fs.AddInitCallback("初始化任务组监听", func() {
		job.InitTaskGroupMonitor()
	})

	fs.AddInitCallback("初始化客户端监听", func() {
		job.InitClientMonitor()
	})

	fs.AddInitCallback("选举", func() {
		// 抢占锁，谁抢到，谁就是master
		go container.Resolve[schedule.Repository]().Election(func() {
			// 推送当前选举结果
			_ = container.Resolve[core.IEvent]("ClusterLeader").Publish(core.AppId)
			<-fs.Context.Done()
		})
	})

	// 每小时检查客户端是否永久离线
	fs.AddInitCallback("20分钟检查客户端是否永久离线", func() {
		tasks.RunNow("RemoveClientJob", 20*time.Minute, job.RemoveClientJob, context.Background())
	})
}
