package job

import (
	"FSchedule/domain/taskGroup"

	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/tasks"
)

// ClearHisTaskJob 自动清除历史任务记录
func ClearHisTaskJob(context *tasks.TaskContext) {
	reservedTaskCount := configure.GetInt("FSchedule.ReservedTaskCount")
	taskGroupRepository := container.Resolve[taskGroup.Repository]()

	// 获取每个任务组，reservedTaskCount之后的TaskId，用于删除
	groupFinishTask, err := taskGroupRepository.GetLastFinishTaskId(reservedTaskCount)
	exception.ThrowRefuseExceptionError(err)
	for taskGroupName, taskId := range groupFinishTask {
		// 清除历史记录
		taskGroupRepository.TaskClearFinish(taskGroupName, taskId)
	}
}
