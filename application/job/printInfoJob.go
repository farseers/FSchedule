package job

import (
	"FSchedule/domain"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
)

// PrintInfoJob 打印客户端、任务组信息
func PrintInfoJob(context *tasks.TaskContext) {
	flog.Infof("当前连接的客户端有 %d 个，监控的任务组 %d 个", domain.ClientCount(), domain.TaskGroupCount())
}
