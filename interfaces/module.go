package interfaces

import (
	"FSchedule/application"
	"FSchedule/interfaces/job"
	"context"
	"time"

	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/tasks"
	"github.com/farseer-go/webapi"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{webapi.Module{}, application.Module{}}
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
}
