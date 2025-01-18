package interfaces

import (
	"FSchedule/application"
	"FSchedule/domain/schedule"
	"FSchedule/interfaces/job"
	"context"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/monitor"
	"github.com/farseer-go/tasks"
	"github.com/farseer-go/webapi"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{webapi.Module{}, application.Module{}, monitor.Module{}}
}

func (module Module) PostInitialize() {
	// // 打印客户端、任务组信息
	// fs.AddInitCallback("启动 【打印客户端、任务组信息】任务", func() {
	// 	tasks.Run("PrintInfoJob", 10*time.Second, job.PrintInfoJob, context.Background())
	// })

	// 10秒更新一次服务端信息
	fs.AddInitCallback("启动 【每10秒更新节点活跃时间】任务", func() {
		tasks.Run("ServerNodeJob", 10*time.Second, job.ServerActivateJob, context.Background())
	})

	// 抢占锁，谁抢到，谁负责这个任务组监控（只允许一个集群节点监控任务组）
	container.Resolve[schedule.Repository]().Schedule("TaskGroupMonitor", func() {
		flog.Infof("开启监控任务组超时监测")
		// 监控任务组超时
		monitor.AddMonitor(1*time.Minute, func() collections.Dictionary[string, any] {
			return job.TaskGroupMonitor()
		})
	})
}
