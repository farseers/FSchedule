package application

import (
	"FSchedule/application/job"
	"FSchedule/domain"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
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

	// 计算任务组的平均耗时
	tasks.Run("SyncAvgSpeedJob", 30*time.Minute, job.SyncAvgSpeedJob, fs.Context)

	// 自动清除历史任务记录
	if configure.GetInt("FSchedule.ReservedTaskCount") > 0 {
		tasks.Run("ClearHisTaskJob", 1*time.Hour, job.ClearHisTaskJob, fs.Context)
	}

	// 打印客户端、任务组信息
	tasks.Run("PrintInfoJob", 10*time.Second, job.PrintInfoJob, fs.Context)

	// 10秒更新一次服务端信息
	tasks.Run("ServerNodeJob", 10*time.Second, job.ServerNodeJob, fs.Context)

	fs.AddInitCallback("初始化任务组监听", func() {
		job.InitTaskGroupMonitor()
	})

	fs.AddInitCallback("初始化客户端监听", func() {
		job.InitClientMonitor()
	})
}

func (module Module) Shutdown() {
}
