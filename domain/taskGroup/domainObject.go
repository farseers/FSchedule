package taskGroup

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/snowflake"
	"github.com/robfig/cron/v3"
	"time"
)

var standardParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type DomainObject struct {
	Id          int64                                  // 主键ID
	Name        string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver         int                                    // 版本
	Task        TaskEO                                 // 最新的任务
	Caption     string                                 // 任务组标题
	Data        collections.Dictionary[string, string] // 本次执行任务时的Data数据
	StartAt     time.Time                              // 开始时间
	NextAt      time.Time                              // 下次执行时间
	Cron        string                                 // 时间定时器表达式
	ActivateAt  time.Time                              // 活动时间
	LastRunAt   time.Time                              // 最后一次完成时间
	IsEnable    bool                                   // 是否开启
	RunSpeedAvg int64                                  // 运行平均耗时
	RunCount    int                                    // 运行次数
	NeedSave    bool                                   // 是否需要保存
}

// UpdateVer 更新新的版本
func (receiver *DomainObject) UpdateVer(name string, caption string, ver int, strCron string, StartAt int64, enable bool) {
	// 只更新高一个版本号的数据
	if receiver.Ver+1 == ver {
		receiver.Name = name
		receiver.Caption = caption
		receiver.Ver = ver
		receiver.Cron = strCron
		receiver.StartAt = time.Unix(StartAt, 0)
		receiver.NeedSave = true
		receiver.IsEnable = enable

		if enable {
			cornSchedule, err := standardParser.Parse(receiver.Cron)
			if err != nil {
				_ = flog.Errorf("Name:%s，Cron格式错误:%s", receiver.Name, receiver.Cron)
				receiver.NeedSave = false
				return
			} else {
				receiver.NextAt = cornSchedule.Next(time.Now())
				receiver.ActivateAt = time.Now()
				receiver.LastRunAt = time.Now()
			}
		}
	}

	if enable && receiver.Task.IsNull() {
		receiver.CreateTask()
		receiver.NeedSave = true
	}
}

// CreateTask 创建新的Task
func (receiver *DomainObject) CreateTask() {
	if receiver.Task.IsFinish() {
		receiver.RunCount++
		receiver.LastRunAt = time.Now()
		receiver.ActivateAt = time.Now()
	}
	receiver.Task = TaskEO{
		Id:          snowflake.GenerateId(),
		Ver:         receiver.Ver,
		Caption:     receiver.Caption,
		Name:        receiver.Name,
		TaskGroupId: receiver.Id,
		StartAt:     receiver.NextAt,
		RunAt:       time.Now(),
		RunSpeed:    0,
		Progress:    0,
		Status:      enum.None,
		CreateAt:    time.Now(),
		SchedulerAt: time.Now(),
		Data:        receiver.Data,
	}
}

// SetClient 分配客户端
func (receiver *DomainObject) SetClient(client ClientVO) {
	receiver.Task.Client = client
	receiver.Task.Status = enum.Working
	receiver.Task.SchedulerAt = time.Now()
	receiver.Task.RunAt = time.Now()
}

// IsNil 不存在
func (receiver *DomainObject) IsNil() bool {
	return receiver.Name == ""
}

// UpdateTask 更新任务信息
func (receiver *DomainObject) UpdateTask(taskEO TaskEO) {
	if taskEO.Id >= receiver.Task.Id {
		receiver.Data = taskEO.Data
		receiver.Task = taskEO
	}
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
	return !receiver.Task.IsNull() &&
		(receiver.Task.Status == enum.None ||
			receiver.Task.Status == enum.ScheduleFail ||
			receiver.Task.Status == enum.Scheduling) &&
		receiver.IsEnable &&
		time.Now().After(receiver.StartAt)
}

// CalculateNextAtByUnix 重新计算下一个执行周期
func (receiver *DomainObject) CalculateNextAtByUnix(timespan int64) {
	if timespan > 0 {
		receiver.NextAt = time.UnixMilli(timespan)
	}
}

// CalculateNextAtByCron 重新计算下一个执行周期
func (receiver *DomainObject) CalculateNextAtByCron() {
	if time.Now().After(receiver.NextAt) {
		cornSchedule, err := standardParser.Parse(receiver.Cron)
		if err != nil {
			_ = flog.Errorf("Name:%s，Cron格式错误:%s", receiver.Name, receiver.Cron)
		}
		receiver.NextAt = cornSchedule.Next(time.Now())
	}
}

// SyncData 同步Data
func (receiver *DomainObject) SyncData() {
	if receiver.Task.Status == enum.Success {
		receiver.Data = receiver.Task.Data
	}
}

// Report 任务报告
func (receiver *DomainObject) Report(status enum.TaskStatus, data collections.Dictionary[string, string], progress int, runSpeed int64, nextTimespan int64, taskGroupRepository Repository) {
	receiver.ActivateAt = time.Now()
	receiver.LastRunAt = time.Now()
	receiver.Task.UpdateTask(status, data, progress, runSpeed)
	receiver.SyncData()
	// 客户端动态计算下一个执行周期
	receiver.CalculateNextAtByUnix(nextTimespan)
	taskGroupRepository.Save(*receiver)

	//if receiver.Task.IsFinish() {
	//	_ = container.Resolve[core.IEvent]("TaskFinish").Publish(receiver)
	//} else {
	//	taskGroupRepository.Save(*receiver)
	//}
}
