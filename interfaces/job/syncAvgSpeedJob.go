package job

import (
	"FSchedule/domain/taskGroup"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/tasks"
)

// SyncAvgSpeedJob 计算任务组的平均耗时
func SyncAvgSpeedJob(context *tasks.TaskContext) {
	repository := container.Resolve[taskGroup.Repository]()
	var speedList = repository.ToTaskSpeedList()
	speedList.Foreach(func(item *taskGroup.TaskEO) {
		if item.RunSpeed > 0 {
			var do = repository.ToEntity(item.Name)
			if !do.IsNil() {
				do.RunSpeedAvg = item.RunSpeed
				repository.Save(do)

				// 发到所有节点上
				_ = container.Resolve[core.IEvent]("TaskGroupUpdate").Publish(do)
			}
		}
	})
}
