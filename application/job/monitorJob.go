package job

import (
	"FSchedule/domain/monitor"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
)

// 加入到监控的列表
var taskGroupList = collections.NewDictionary[string, chan monitor.DomainObject]()

// MonitorJob 任务组监听
func MonitorJob() {
	repository := container.Resolve[taskGroup.Repository]()
	lst := repository.ToList()
	for _, taskGroupDO := range lst.ToArray() {
		scheduleDO := mapper.Single[monitor.DomainObject](taskGroupDO)

		// 新的任务组不再当前列表，说明被其它节点处理了。
		if !taskGroupList.ContainsKey(taskGroupDO.Name) {
			c := make(chan monitor.DomainObject)
			taskGroupList.Add(taskGroupDO.Name, c)
			flog.Infof("任务组：%s %s 加入调度线程", taskGroupDO.Name, taskGroupDO.Caption)

			go container.ResolveIns(&scheduleDO).Start(c)
		}
		// 将最新的任务组数据发送到通道
		c := taskGroupList.GetValue(taskGroupDO.Name)
		c <- scheduleDO
	}
}
