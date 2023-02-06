package domainEvent

import (
	"FSchedule/domain"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
)

func TaskFinishEvent(message any, _ eventBus.EventArgs) {
	do := message.(*domain.TaskGroupMonitor)
	// 加锁
	scheduleRepository := container.Resolve[schedule.Repository]()
	lock := scheduleRepository.GetLock(do.Name)
	if !lock.TryLock() {
		flog.Debugf("创建新任务时加锁失败，Job=%s，serverIP=%s，serverId=%d", do.Name, fs.AppIp, fs.AppId)
		return
	}
	defer lock.ReleaseLock()

	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	taskGroupDO := taskGroupRepository.ToEntity(do.Name)
	// 任务初始化
	taskGroupDO.CreateTask()
	taskGroupRepository.SaveAndTask(taskGroupDO)
}
