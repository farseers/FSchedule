package job

import (
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

// SyncAvgSpeedJob 计算任务组的平均耗时
func SyncAvgSpeedJob(context *tasks.TaskContext) {
	repository := container.Resolve[taskGroup.Repository]()
	for _, taskGroupDO := range repository.ToList().ToArray() {
		var speedList = repository.ToTaskSpeedList(taskGroupDO.Name)
		var runSpeedAvg = taskGroup.NewTaskSpeed(speedList).GetAvgSpeed()

		if runSpeedAvg > 0 {
			var do = repository.ToEntity(taskGroupDO.Name)
			if !do.IsNil() {
				do.RunSpeedAvg = runSpeedAvg
				repository.Save(do)
			}
		}
	}
}
