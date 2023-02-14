package job

import (
	"FSchedule/domain"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
)

// PrintInfoJob 打印客户端、任务组信息
func PrintInfoJob(context *tasks.TaskContext) {
	flog.Infof("当前%d个客户端连接，%d个任务组在监控，（%d个在执行）", domain.ClientCount(), domain.TaskGroupCount(), domain.TaskGroupEnableCount())
}
