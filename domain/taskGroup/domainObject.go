package taskGroup

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/snowflake"
	"github.com/robfig/cron/v3"
	"time"
)

var StandardParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type DomainObject struct {
	Id          int64                                  // 主键ID
	Name        string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver         int                                    // 版本
	Task        TaskEO                                 // 最新的任务
	Caption     string                                 // 任务组标题
	Data        collections.Dictionary[string, string] // 本次执行任务时的Data数据
	StartAt     dateTime.DateTime                      // 开始时间
	NextAt      dateTime.DateTime                      // 下次执行时间
	Cron        string                                 // 时间定时器表达式
	ActivateAt  dateTime.DateTime                      // 活动时间
	LastRunAt   dateTime.DateTime                      // 最后一次完成时间
	IsEnable    bool                                   // 是否开启
	RunSpeedAvg int64                                  // 运行平均耗时
	RunCount    int                                    // 运行次数
	NeedSave    bool                                   // 是否需要保存（API接口使用）
}

func New(name string, caption string, ver int, strCron string, data collections.Dictionary[string, string], startAt int64, enable bool) *DomainObject {
	do := &DomainObject{}
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
				_ = flog.Errorf("任务组:%s（%d），Cron格式错误:%s", receiver.Name, receiver.Id, receiver.Cron)
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
}

func (receiver *DomainObject) Update() {
	// 已下发的任务，不能修改
	switch receiver.Task.Status {
	case enum.None, enum.Fail, enum.Success, enum.ScheduleFail:
		cornSchedule, err := StandardParser.Parse(receiver.Cron)
		if err != nil {
			_ = flog.Errorf("任务组:%s（%d），Cron格式错误:%s", receiver.Name, receiver.Id, receiver.Cron)
		}
		receiver.NextAt = dateTime.New(cornSchedule.Next(time.Now()))
		receiver.Task.Data = receiver.Data
	}
}

// CreateTask 创建新的Task
func (receiver *DomainObject) CreateTask() {
	if receiver.Task.IsFinish() {
		receiver.RunCount++
		receiver.LastRunAt = dateTime.Now()
		receiver.ActivateAt = dateTime.Now()
	}
	receiver.Task = TaskEO{
		Id:          snowflake.GenerateId(),
		Ver:         receiver.Ver,
		Caption:     receiver.Caption,
		Name:        receiver.Name,
		TaskGroupId: receiver.Id,
		StartAt:     receiver.NextAt,
		RunAt:       dateTime.Now(),
		RunSpeed:    0,
		Progress:    0,
		Status:      enum.None,
		CreateAt:    dateTime.Now(),
		SchedulerAt: dateTime.Now(),
		Data:        receiver.Data,
	}
}

// SetClient 分配客户端
func (receiver *DomainObject) SetClient(client ClientVO) {
	receiver.Task.Client = client
	receiver.Task.Status = enum.Working
	receiver.Task.SchedulerAt = dateTime.Now()
	receiver.Task.RunAt = dateTime.Now()
	// 重新赋值是为了担心数据被手动改了
	receiver.Task.Data = receiver.Data
	receiver.Task.Name = receiver.Name
	receiver.Task.TaskGroupId = receiver.Id
}

// IsNil 不存在
func (receiver *DomainObject) IsNil() bool {
	return receiver.Name == ""
}

// ScheduleFail 调度失败
func (receiver *DomainObject) ScheduleFail() {
	receiver.Task.ScheduleFail()
}

// ClientOffline 客户端下线了
func (receiver *DomainObject) ClientOffline() {
	receiver.Task.SetFail()
}

// CanScheduler 是否可以调度
func (receiver *DomainObject) CanScheduler() bool {
	now := dateTime.Now()
	return !receiver.Task.IsNull() &&
		(receiver.Task.Status == enum.None ||
			receiver.Task.Status == enum.ScheduleFail ||
			receiver.Task.Status == enum.Scheduling) &&
		receiver.IsEnable &&
		now.After(receiver.StartAt)
}

// CalculateNextAtByUnix 重新计算下一个执行周期
func (receiver *DomainObject) CalculateNextAtByUnix(timespan int64) {
	if timespan > 0 {
		receiver.NextAt = dateTime.NewUnixMilli(timespan)
	}
}

// CalculateNextAtByCron 重新计算下一个执行周期
func (receiver *DomainObject) CalculateNextAtByCron() {
	now := dateTime.Now()
	// 成功才要计算下一个周期
	if now.After(receiver.NextAt) {
		switch receiver.Task.Status {
		case enum.Success:
			cornSchedule, err := StandardParser.Parse(receiver.Cron)
			if err != nil {
				_ = flog.Errorf("任务组:%s（%d），Cron格式错误:%s", receiver.Name, receiver.Id, receiver.Cron)
			}
			receiver.NextAt = dateTime.New(cornSchedule.Next(time.Now()))
		case enum.Fail:
			// 失败，则为下一秒在执行
			receiver.NextAt = now.AddSeconds(1)
		}
	}
}

// SyncData 同步Data
func (receiver *DomainObject) SyncData() {
	if receiver.Task.Status == enum.Success || receiver.Task.Status == enum.Fail {
		receiver.Data = receiver.Task.Data
		if receiver.Data.Count() == 0 {
			flog.Infof("")
		}
	}
}

// Report 任务报告
func (receiver *DomainObject) Report(status enum.TaskStatus, data collections.Dictionary[string, string], progress int, runSpeed int64, nextTimespan int64, taskGroupRepository Repository) {
	receiver.ActivateAt = dateTime.Now()
	receiver.LastRunAt = dateTime.Now()
	receiver.Task.UpdateTask(status, data, progress, runSpeed)
	receiver.SyncData()
	// 客户端动态计算下一个执行周期
	receiver.CalculateNextAtByUnix(nextTimespan)
	taskGroupRepository.Save(*receiver)
}

// ReportFail 任务报告，未找到任务
func (receiver *DomainObject) ReportFail(taskGroupRepository Repository) {
	receiver.Task.UpdateTask(enum.Fail, receiver.Task.Data, receiver.Task.Progress, receiver.Task.RunSpeed)
	taskGroupRepository.Save(*receiver)
}
