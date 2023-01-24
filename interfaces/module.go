package interfaces

import (
	"FSchedule/application"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/webapi"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{webapi.Module{}, application.Module{}}
}

func (module Module) PreInitialize() {
	//TODO implement me
}

func (module Module) Initialize() {
	//TODO implement me
}

func (module Module) PostInitialize() {
	//TODO implement me
}

func (module Module) Shutdown() {
	//TODO implement me
}
