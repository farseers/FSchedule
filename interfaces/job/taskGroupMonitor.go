package job

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/taskGroup"
	"fmt"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
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

	lstStr := collections.NewList[string]()
	lst.Foreach(func(item *taskGroup.DomainObject) {
		switch item.Task.ExecuteStatus {
		// 超时未执行
		case executeStatus.None:
			if item.Task.StartAt.Before(dateTime.Now()) {
				lstStr.Add(fmt.Sprintf("%s(%s)\r\n超时%s未执行。\r\n", item.Caption, item.Name, (time.Duration(dateTime.Now().Sub(item.Task.StartAt).Seconds()) * time.Second).String()))
			}
		// 执行超时
		case executeStatus.Working:
			executeTime := (time.Duration(dateTime.Now().Sub(item.Task.SchedulerAt).Seconds()) * time.Second)
			if float64(executeTime.Milliseconds())*1.3 > float64(item.RunSpeedAvg) {
				lstStr.Add(fmt.Sprintf("%s(%s)\r\n执行了%s。\r\n", item.Caption, item.Name, executeTime.String()))
			}

		}
	})

	if lstStr.Count() > 0 {
		dic.Add("fschedule_timeout", lstStr.ToString("\r\n"))
	}
	return dic
}
