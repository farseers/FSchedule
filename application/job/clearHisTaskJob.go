package job

import (
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

// ClearHisTaskJob 自动清除历史任务记录
func ClearHisTaskJob(context *tasks.TaskContext) {
	reservedTaskCount := configure.GetInt("FSchedule.ReservedTaskCount")
	taskGroupRepository := container.Resolve[taskGroup.Repository]()

	//curIndex := 0
	result := 0
	taskGroupRepository.ToList().Foreach(func(taskGroupDO *taskGroup.DomainObject) {
		lstTask := taskGroupRepository.ToFinishList(taskGroupDO.Id, reservedTaskCount)
		if lstTask.Count() == 0 {
			return
		}

		result += lstTask.Count()
		var taskId = lstTask.Min(func(item taskGroup.TaskEO) any {
			return item.Id
		}).(int)

		// 清除历史记录
		taskGroupRepository.ClearFinish(taskGroupDO.Id, taskId)
	})
}
