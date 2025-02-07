package job

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/taskGroup"
	"fmt"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
)

// TaskGroupMonitor 监控任务组超时
func TaskGroupMonitor() collections.Dictionary[string, any] {
	// 发送消息
	dic := collections.NewDictionary[string, any]()

	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	lst := taskGroupRepository.ToList()
	// 状态为可用、非完成状态、并按开始时间排序
	lst = lst.Where(func(item taskGroup.DomainObject) bool {
		return item.IsEnable && !item.Task.IsFinish()
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.Name
	}).OrderBy(func(item taskGroup.DomainObject) any {
		return item.Task.StartAt.UnixNano()
	}).ToList()

	lstUnWork := collections.NewList[string]()
	lstTimeout := collections.NewList[string]()
	lst.Foreach(func(item *taskGroup.DomainObject) {
		switch item.Task.ExecuteStatus {
		// 超时未执行
		case executeStatus.None:
			if difference := dateTime.Now().Sub(item.Task.StartAt); difference.Seconds() > 5 {
				lstUnWork.Add(fmt.Sprintf("%s(%s)\r\n超时%s未执行。\r\n", item.Caption, item.Name, difference.String()))
				// 发到所有节点上，主动通知到任务组监控，用于激活任务
				_ = container.Resolve[core.IEvent]("TaskGroupUpdate").Publish(*item)
			}
		// 执行超时
		case executeStatus.Working:
			// 执行的时间
			executeTime := dateTime.Now().Sub(item.Task.SchedulerAt)
			// 比平均时间多1.5倍，且大于1分钟，则告警
			if difference := float64(executeTime.Milliseconds()) - float64(item.RunSpeedAvg)*1.5; difference > 60000 {
				lstTimeout.Add(fmt.Sprintf("%s(%s)\r\n执行了%s。\r\n", item.Caption, item.Name, executeTime.String()))
			}
		}
	})

	if lstUnWork.Count() > 0 {
		dic.Add("fschedule_unwork", lstUnWork.ToString("\r\n"))
	}

	if lstTimeout.Count() > 0 {
		dic.Add("fschedule_timeout", lstTimeout.ToString("\r\n"))
	}
	return dic
}
