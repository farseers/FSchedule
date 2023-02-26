package job

import (
	"FSchedule/domain"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
)

// PrintInfoJob 打印客户端、任务组信息
func PrintInfoJob(context *tasks.TaskContext) {
	flog.Infof("当前%s个客户端连接（%s个正常状态），%s个任务组在监控（%s个在执行）", flog.Red(domain.ClientCount()), flog.Green(domain.ClientNormalCount()), flog.Red(domain.TaskGroupCount()), flog.Green(domain.TaskGroupEnableCount()))
}
