package job

import (
	"FSchedule/domain/serverNode"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

// ServerActivateJob 更新节点信息
func ServerActivateJob(context *tasks.TaskContext) {
	repository := container.Resolve[serverNode.Repository]()
	serverNodeDO := repository.ToEntity(fs.AppId)
	serverNodeDO.Activate()
	repository.Save(&serverNodeDO)
}
