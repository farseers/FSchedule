package infrastructure

import (
	"FSchedule/domain/serverNode"
	"FSchedule/infrastructure/http"
	"FSchedule/infrastructure/localQueue"
	"FSchedule/infrastructure/repository"
	"FSchedule/infrastructure/repository/context"
	"FSchedule/interfaces/domainEvent"
	"github.com/farseer-go/data"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fs/timingWheel"
	"github.com/farseer-go/linkTrace"
	"github.com/farseer-go/queue"
	"github.com/farseer-go/redis"
	"time"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{data.Module{}, redis.Module{}, eventBus.Module{}, queue.Module{}, linkTrace.Module{}}
}

func (module Module) PostInitialize() {
	timingWheel.Start()

	// 注册上下文
	context.InitMysqlContext()
	context.InitRedisContext()

	// 注册仓储
	repository.InitRepository()

	// 注册选举事件
	redis.RegisterEvent("default", "ClusterLeader").RegisterSubscribe("选举事件", domainEvent.ClusterLeaderSubscribe)

	// 注册任务组更新通知事件
	redis.RegisterEvent("default", "TaskGroupUpdate").RegisterSubscribe("任务组有更新", domainEvent.TaskGroupUpdateSubscribe)

	// 队列任务日志
	queue.Subscribe("TaskLogQueue", "同步日志到数据库", 1000, 5*time.Second, localQueue.TaskLogQueueConsumer)

	// 注册客户端http
	http.InitHttp()

	fs.AddInitCallback("注册节点信息", func() {
		container.Resolve[serverNode.Repository]().Save(serverNode.New())
	})
}
