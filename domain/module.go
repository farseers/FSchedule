package domain

import (
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fs/timingWheel"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return nil
}

func (module Module) PreInitialize() {
	timingWheel.Start()
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
}

func (module Module) Shutdown() {
}
