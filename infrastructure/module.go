package infrastructure

import (
	"FSchedule/application/domainEvent"
	"FSchedule/application/job"
	"FSchedule/infrastructure/http"
	"FSchedule/infrastructure/repository"
	"github.com/farseer-go/data"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/queue"
	"github.com/farseer-go/redis"
	"github.com/farseer-go/tasks"
	"time"
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

	// 检查超时离线的客户端
	tasks.Run("CheckClientOffline", 3*time.Second, job.MonitorClientJob, fs.Context)
	// 任务组监听
	tasks.Run("MonitorTaskGroup", 3*time.Second, job.MonitorTaskGroupJob, fs.Context)

	// 客户端离线通知
	eventBus.RegisterEvent("ClientOffline", domainEvent.RemoveClientEvent)
	// 任务状态有变更
	eventBus.RegisterEvent("TaskScheduler", domainEvent.SchedulerEvent)
	// 检查进行中的任务
	eventBus.RegisterEvent("CheckWorking", domainEvent.CheckWorkingEvent)
	// 任务完成事件
	eventBus.RegisterEvent("TaskFinish", domainEvent.TaskFinishEvent)

	// 注册客户端更新通知事件
	redis.RegisterEvent("default", "ClientUpdate", domainEvent.ClientUpdateEvent)
	// 注册任务组更新通知事件
	redis.RegisterEvent("default", "TaskGroupUpdate", domainEvent.TaskGroupUpdateEvent)
}

func (module Module) Shutdown() {
}
