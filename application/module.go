package application

import (
	"FSchedule/application/job"
	"FSchedule/domain"
	"github.com/farseer-go/fs/modules"
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
	job.InitTaskGroupMonitor()
	job.InitClientMonitor()
}

func (module Module) Shutdown() {
}
