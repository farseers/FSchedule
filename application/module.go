package application

import (
	"FSchedule/domain"
	"FSchedule/domain/schedule"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{domain.Module{}}
}

func (module Module) PostInitialize() {
	fs.AddInitCallback("选举", func() {
		// 抢占锁，谁抢到，谁就是master
		go container.Resolve[schedule.Repository]().Election(func() {
			// 推送当前选举结果
			_ = container.Resolve[core.IEvent]("ClusterLeader").Publish(core.AppId)
			<-fs.Context.Done()
		})
	})
}
