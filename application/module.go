package application

import (
	"FSchedule/application/job"
	"FSchedule/domain"
	"github.com/farseer-go/fs"
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
	tasks.Run("CheckClientOffline", 3*time.Second, job.CheckClientOfflineJob, fs.Context)
}

func (module Module) Shutdown() {
}
