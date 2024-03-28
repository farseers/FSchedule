package taskGroup

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/snowflake"
	"github.com/robfig/cron/v3"
	"time"
)

var StandardParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type DomainObject struct {
	Name              string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver               int                                    // 版本
	Task              TaskEO                                 // 最新的任务
	Caption           string                                 // 任务组标题
	Data              collections.Dictionary[string, string] // 本次执行任务时的Data数据
	StartAt           dateTime.DateTime                      // 开始时间
	NextAt            dateTime.DateTime                      // 下次执行时间
	Cron              string                                 // 时间定时器表达式
	ActivateAt        dateTime.DateTime                      // 活动时间
	LastRunAt         dateTime.DateTime                      // 最后一次完成时间
	LastExecuteStatus executeStatus.Enum                     // 上次执行结果
	IsEnable          bool                                   // 是否开启
	RunSpeedAvg       int64                                  // 运行平均耗时
	RunCount          int                                    // 运行次数
	NeedSave          bool                                   // 是否需要保存（API接口使用）
	RetryDelaySecond  int                                    // 失败后多少秒重试（0不重试）
}

func New(name string, caption string, ver int, strCron string, data collections.Dictionary[string, string], startAt int64, enable bool) *DomainObject {
	do := &DomainObject{Ver: ver - 1}
	do.UpdateVer(name, caption, ver, strCron, data, startAt, enable)
	return do
}

// UpdateVer 更新新的版本
func (receiver *DomainObject) UpdateVer(name string, caption string, ver int, strCron string, data collections.Dictionary[string, string], startAt int64, enable bool) {
	// 只更新高一个版本号的数据
	if receiver.Ver+1 == ver {
		receiver.Name = name
		receiver.Caption = caption
		receiver.Ver = ver
		receiver.Cron = strCron
		receiver.StartAt = dateTime.NewUnix(startAt)
		receiver.NeedSave = true
		receiver.IsEnable = enable
		receiver.Data = data

		if enable {
			cornSchedule, err := StandardParser.Parse(receiver.Cron)
			if err != nil {
				_ = flog.Errorf("任务组:%s，Cron格式错误:%s", receiver.Name, receiver.Cron)
				receiver.NeedSave = false
				return
			} else {
				receiver.NextAt = dateTime.New(cornSchedule.Next(time.Now()))
				receiver.ActivateAt = dateTime.Now()
				receiver.LastRunAt = dateTime.Now()
			}
		}
	}

	if enable && receiver.Task.IsNull() {
		receiver.CreateTask()
		receiver.NeedSave = true
	}

	if receiver.StartAt.Year() < 2000 {
		receiver.StartAt = dateTime.Now()
	}

	if receiver.ActivateAt.Year() < 2000 {
		receiver.ActivateAt = dateTime.Now()
	}

	if receiver.LastRunAt.Year() < 2000 {
		receiver.LastRunAt = dateTime.Now()
	}

	if receiver.NextAt.Year() < 2000 {
		receiver.NextAt = dateTime.Now()
	}
}

// 调度中 或 执行中
func (receiver *DomainObject) IsScheduleOrWorking() bool {
	return (receiver.Task.ScheduleStatus == scheduleStatus.Scheduling || receiver.Task.ScheduleStatus == scheduleStatus.Success) &&
		(receiver.Task.ExecuteStatus == executeStatus.Working || receiver.Task.ExecuteStatus == executeStatus.None)
}

// 更新时间
func (receiver *DomainObject) Update() {
	// 调度中 或 执行中 不允许修改
	if receiver.IsScheduleOrWorking() {
		exception.ThrowWebException(403, "任务组处于调度状态或执行中，不允许修改")
	}

	cornSchedule, err := StandardParser.Parse(receiver.Cron)
	if err != nil {
		_ = flog.Errorf("任务组:%s，Cron格式错误:%s", receiver.Name, receiver.Cron)
	}
	receiver.NextAt = dateTime.New(cornSchedule.Next(time.Now()))
	receiver.Task.Data = receiver.Data
	receiver.Task.StartAt = receiver.NextAt
}

// CreateTask 创建新的Task
func (receiver *DomainObject) CreateTask() {
	if receiver.Task.IsFinish() {
		receiver.RunCount++
		receiver.LastRunAt = dateTime.Now()
		receiver.ActivateAt = dateTime.Now()
		receiver.LastExecuteStatus = receiver.Task.ExecuteStatus
	}
	receiver.Task = TaskEO{
		Id:             snowflake.GenerateId(),
		Ver:            receiver.Ver,
		Caption:        receiver.Caption,
		Name:           receiver.Name,
		StartAt:        receiver.NextAt,
		RunAt:          dateTime.Now(),
		RunSpeed:       0,
		Progress:       0,
		ScheduleStatus: scheduleStatus.None,
		ExecuteStatus:  executeStatus.None,
		CreateAt:       dateTime.Now(),
		SchedulerAt:    dateTime.Now(),
		Data:           receiver.Data,
	}
}

// IsNil 不存在
func (receiver *DomainObject) IsNil() bool {
	return receiver.Name == ""
}

// CanScheduler 是否可以调度
func (receiver *DomainObject) CanScheduler() bool {
	now := dateTime.Now()
	return !receiver.Task.IsNull() &&
		(receiver.Task.ScheduleStatus == scheduleStatus.Scheduling || receiver.Task.ScheduleStatus == scheduleStatus.None) &&
		receiver.Task.ExecuteStatus == executeStatus.None &&
		receiver.IsEnable &&
		now.After(receiver.StartAt)
}

// CalculateNextAtByUnix 重新计算下一个执行周期
func (receiver *DomainObject) CalculateNextAtByUnix(timespan int64) {
	if timespan > 0 {
		flog.Infof("任务组:%s 设置了timespan=%d,更新下一次时间为%s", receiver.Name, timespan, receiver.NextAt.ToString("yyyy-MM-dd HH:mm:ss"))
		receiver.NextAt = dateTime.NewUnixMilli(timespan)
	}
}

// CalculateNextAtByCron 重新计算下一个执行周期
func (receiver *DomainObject) CalculateNextAtByCron() bool {
	// 时间相等，说明客户端没有设置过时间
	if receiver.NeverSetNextAt() {
		// 执行结果为失败，且设置了重试。则按重试时间计算下次执行时间。
		if receiver.Task.ExecuteStatus == executeStatus.Fail && receiver.RetryDelaySecond > 0 {
			// 失败，则为下一秒在执行
			receiver.NextAt = dateTime.Now().AddSeconds(receiver.RetryDelaySecond)
			flog.Infof("任务组:%s 执行失败，将在1秒后（%s）继续执行", receiver.Name, receiver.NextAt.ToString("yyyy-MM-dd HH:mm:ss"))
			return true
		}

		cornSchedule, err := StandardParser.Parse(receiver.Cron)
		if err != nil {
			_ = flog.Errorf("任务组:%s，Cron格式错误:%s，已将任务暂停。", receiver.Name, receiver.Cron)
			receiver.IsEnable = false
			return false
		} else {
			receiver.NextAt = dateTime.New(cornSchedule.Next(time.Now()))
		}
	}
	return true
}

// SyncData 同步Data
func (receiver *DomainObject) SyncData() {
	if receiver.Task.IsFinish() {
		receiver.Data = receiver.Task.Data
	}
}

// Report 任务报告
func (receiver *DomainObject) Report(status executeStatus.Enum, data collections.Dictionary[string, string], progress int, runSpeed int64, nextTimespan int64, remark string, taskGroupRepository Repository) {
	receiver.ActivateAt = dateTime.Now()
	receiver.LastRunAt = dateTime.Now()
	receiver.Task.UpdateTask(status, data, progress, runSpeed, remark)
	receiver.SyncData()
	// 客户端动态计算下一个执行周期
	receiver.CalculateNextAtByUnix(nextTimespan)
	taskGroupRepository.Save(*receiver)
}

// ReportFail 任务报告，未找到任务
func (receiver *DomainObject) ReportFail(taskGroupRepository Repository, remark string) {
	receiver.Task.UpdateTaskStatus(executeStatus.Fail, remark)
	taskGroupRepository.Save(*receiver)
}

// 设置状态
func (receiver *DomainObject) SetEnable(enable bool) {
	receiver.IsEnable = enable
}

// 没有设置过时间
func (receiver *DomainObject) NeverSetNextAt() bool {
	return receiver.NextAt.UnixMilli() == receiver.Task.StartAt.UnixMilli()
}
