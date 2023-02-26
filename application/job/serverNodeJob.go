package job

import (
	"FSchedule/domain/serverNode"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

// ServerNodeJob 注册服务端信息
func ServerNodeJob(context *tasks.TaskContext) {
	do := serverNode.New()
	container.Resolve[serverNode.Repository]().Save(do)
}
