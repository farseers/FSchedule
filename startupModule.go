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
	//parse, _ := taskGroup.StandardParser.Parse("* 0/10 * * * ?")
	//n := time.Now()
	//fmt.Println(n)
	//fmt.Println("--------------")
	//for i := 0; i < 10; i++ {
	//	n = parse.Next(n)
	//	fmt.Println(n)
	//}
	//fmt.Println("--------------")
}

func (module StartupModule) Shutdown() {
}
