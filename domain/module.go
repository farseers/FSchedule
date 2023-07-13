package domain

import (
	"github.com/farseer-go/fs/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return nil
}
