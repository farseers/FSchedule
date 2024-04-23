package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/enum/scheduleStatus"
	"FSchedule/domain/taskGroup"
	"fmt"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
	"time"
)

// SchedulerEvent 任务调度
func SchedulerEvent(message any, _ core.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	// 只订阅调度状态的事件
	if do.Task.ScheduleStatus != scheduleStatus.Scheduling {
		return
	}
	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	clientRepository := container.Resolve[client.Repository]()

	for {
		if !do.CanScheduler() {
			flog.Debugf("任务组：%s 条件不满足无法调度，延迟：%s", do.Name, dateTime.Since(do.Task.StartAt).String())
			do.Task.ScheduleFail("条件不满足无法调度")
			return
		}

		// 轮询的方式取到客户端
		clientSchedule := do.PollingClient()
		// 没有可调度的客户端
		if clientSchedule == nil || clientSchedule.IsNil() {
			flog.Debugf("任务组：%s 没有可调度的客户端，延迟：%s", do.Name, dateTime.Since(do.Task.StartAt).String())
			do.Task.ScheduleFail("没有可用的客户端")
			taskGroupRepository.Save(*do.DomainObject)
			return
		}

		// 请求客户端
		clientTask := mapper.Single[client.TaskEO](do.Task)
		if do.Name != clientTask.Name {
			_ = flog.Errorf("任务组：%s 注意，任务调度，发现task.Name不一致，TaskId=%d，taskName=%s, task=%+v", do.Name, clientTask.Id, clientTask.Name, clientTask)
		}
		var err error
		var success bool
		if success, err = clientSchedule.TrySchedule(clientTask); success {
			// 调度成功，分配客户端
			do.Task.ScheduleSuccess(mapper.Single[taskGroup.ClientVO](clientSchedule))
			clientRepository.Save(clientSchedule)
			taskGroupRepository.SaveAndTask(*do.DomainObject)
			return
		}

		// 调度失败
		do.Task.ScheduleFail(fmt.Sprintf("请求客户端%s（%d）：%s:%d失败:%s", clientSchedule.Name, clientSchedule.Id, clientSchedule.Ip, clientSchedule.Port, err.Error()))
		taskGroupRepository.Save(*do.DomainObject)
		clientRepository.Save(clientSchedule)

		time.Sleep(100 * time.Millisecond)
	}
}
