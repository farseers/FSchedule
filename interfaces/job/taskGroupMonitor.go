package job

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/taskGroup"
	"fmt"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
)

// TaskGroupMonitor 监控任务组超时
func TaskGroupMonitor() collections.Dictionary[string, any] {
	// 发送消息
	dic := collections.NewDictionary[string, any]()

	taskGroupRepository := container.Resolve[taskGroup.Repository]()
	lst := taskGroupRepository.ToList()
	// 状态为可用、非完成状态、并按开始时间排序
	lst = lst.Where(func(item taskGroup.DomainObject) bool { return item.IsEnable && !item.Task.IsFinish() }).ToList()

	lstStr := collections.NewList[string]()
	lst.Foreach(func(item *taskGroup.DomainObject) {
		switch item.Task.ExecuteStatus {
		// 超时未执行
		case executeStatus.None:
			if item.Task.StartAt.Before(dateTime.Now()) {
				lstStr.Add(fmt.Sprintf("%s(%s)，超时%s未执行。", item.Caption, item.Caption, (time.Duration(dateTime.Now().Sub(item.Task.StartAt).Seconds()) * time.Second).String()))
			}
		// 执行超时
		case executeStatus.Working:
			executeTime := (time.Duration(dateTime.Now().Sub(item.Task.SchedulerAt).Seconds()) * time.Second)
			if float64(executeTime.Milliseconds())*1.3 > float64(item.RunSpeedAvg) {
				lstStr.Add(fmt.Sprintf("%s(%s)，执行了%s。", item.Caption, item.Caption, executeTime.String()))
			}

		}
	})

	if lstStr.Count() > 0 {
		dic.Add("fschedule_timeout", lstStr.ToString(""))
	}
	flog.Info(lstStr.ToString("\n"))
	return dic
}
