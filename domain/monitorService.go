package domain

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

func MonitorContainsTaskGroup(name string) bool {
	return taskGroupList.ContainsKey(name)
}

func MonitorPush(taskGroupDO taskGroup.DomainObject) {
	scheduleDO := mapper.Single[monitor.DomainObject](taskGroupDO)

	// 新的任务组不再当前列表，说明被其它节点处理了。
	if !MonitorContainsTaskGroup(scheduleDO.Name) {
		addTaskGroup(scheduleDO)
	}

	// 将最新的任务组数据发送到通道
	c := taskGroupList.GetValue(scheduleDO.Name)
	c <- scheduleDO
}

func addTaskGroup(scheduleDO monitor.DomainObject) {
	c := make(chan monitor.DomainObject)
	taskGroupList.Add(scheduleDO.Name, c)
	flog.Infof("任务组：%s %s 加入调度线程", scheduleDO.Name)

	go container.ResolveIns(&scheduleDO).Start(c)
}
