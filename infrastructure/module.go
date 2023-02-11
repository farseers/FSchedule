package infrastructure

import (
	"FSchedule/application/domainEvent"
	"FSchedule/infrastructure/http"
	"FSchedule/infrastructure/repository"
	"github.com/farseer-go/data"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/queue"
	"github.com/farseer-go/redis"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{data.Module{}, redis.Module{}, eventBus.Module{}, queue.Module{}}
}

func (module Module) PreInitialize() {
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
	repository.InitRepository()
	http.InitHttp()

	// 任务状态有变更
	eventBus.RegisterEvent("TaskScheduler", domainEvent.SchedulerEvent)
	// 检查进行中的任务
	eventBus.RegisterEvent("CheckWorking", domainEvent.CheckWorkingEvent)
	// 任务完成事件
	eventBus.RegisterEvent("TaskFinish", domainEvent.TaskFinishEvent)
	// 新的客户端加入
	eventBus.RegisterEvent("ClientJoin", domainEvent.ClientJoinEvent)
	// 客户端下线了
	eventBus.RegisterEvent("ClientOffline", domainEvent.ClientOfflineEvent)

	// 注册客户端更新通知事件
	redis.RegisterEvent("default", "ClientUpdate", domainEvent.ClientUpdateSubscribe)

	// 注册任务组更新通知事件
	redis.RegisterEvent("default", "TaskGroupUpdate", domainEvent.TaskGroupUpdateSubscribe)
}

func (module Module) Shutdown() {
}
