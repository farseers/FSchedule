package application

import (
	"FSchedule/application/job"
	"FSchedule/domain"
	"FSchedule/domain/schedule"
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

func (module Module) PreInitialize() {
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
	// 打印客户端、任务组信息
	fs.AddInitCallback("打印客户端、任务组信息", func() {
		tasks.Run("PrintInfoJob", 10*time.Second, job.PrintInfoJob, fs.Context)
	})

	// 10秒更新一次服务端信息
	fs.AddInitCallback("10秒更新一次服务端信息", func() {
		tasks.Run("ServerNodeJob", 10*time.Second, job.ServerActivateJob, fs.Context)
	})

	fs.AddInitCallback("初始化任务组监听", func() {
		job.InitTaskGroupMonitor()
	})

	fs.AddInitCallback("初始化客户端监听", func() {
		job.InitClientMonitor()
	})

	fs.AddInitCallback("选举", func() {
		// 抢占锁，谁抢到，谁就是master
		container.Resolve[schedule.Repository]().Election(func() {
			// 推送当前选举结果
			_ = container.Resolve[core.IEvent]("ClusterLeader").Publish(fs.AppId)
		})
	})
}

func (module Module) Shutdown() {
}
