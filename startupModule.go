package main

import (
	"FSchedule/infrastructure"
	"FSchedule/interfaces"
	"github.com/farseer-go/fs/modules"
)

type StartupModule struct {
}

func (module StartupModule) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{infrastructure.Module{}, interfaces.Module{}}
}

func (module StartupModule) PreInitialize() {
}

func (module StartupModule) Initialize() {
}

func (module StartupModule) PostInitialize() {

}

func (module StartupModule) Shutdown() {
}
