package taskGroup

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
	"time"
)

// TaskEO 任务记录
type TaskEO struct {
	Id          int64                                  // 主键
	TaskGroupId int64                                  // 任务组ID
	Name        string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver         int                                    // 版本
	Caption     string                                 // 任务组标题
	StartAt     time.Time                              // 开始时间
	RunAt       time.Time                              // 实际执行时间
	RunSpeed    int64                                  // 运行耗时
	Client      ClientVO                               // 客户端
	Progress    int                                    // 进度0-100
	Status      enum.TaskStatus                        // 状态
	SchedulerAt time.Time                              // 调度时间
	Data        collections.Dictionary[string, string] // 本次执行任务时的Data数据
	CreateAt    time.Time                              // 任务创建时间
}

func NewTaskDO() *TaskEO {
	return &TaskEO{}
}

// SetClient 调度时设置客户端
func (receiver *TaskEO) SetClient(client ClientVO) {
	receiver.Status = enum.Working
	receiver.SchedulerAt = time.Now()
	receiver.Client = client
}

// SetJobName 更新了JobName，则要立即更新Task的JobName
func (receiver *TaskEO) SetJobName(name string) {
	receiver.Name = name
}

// SetFail 设备为失败
func (receiver *TaskEO) SetFail() {
	receiver.Status = enum.Fail
}

// Scheduling 调度
func (receiver *TaskEO) Scheduling() {
	receiver.Status = enum.Scheduling
}

// ScheduleFail 调度失败
func (receiver *TaskEO) ScheduleFail() {
	receiver.Status = enum.ScheduleFail
}

// IsNull 未分配
func (receiver *TaskEO) IsNull() bool {
	return receiver.Id == 0 && receiver.Caption == "" && receiver.Name == ""
}

// IsFinish 是否完成
func (receiver *TaskEO) IsFinish() bool {
	return receiver.Status == enum.Success || receiver.Status == enum.Fail
}

// IsWorking 是否为执行中
func (receiver *TaskEO) IsWorking() bool {
	return receiver.Status == enum.Working
}

// UpdateTask 更新任务
func (receiver *TaskEO) UpdateTask(status enum.TaskStatus, data collections.Dictionary[string, string], progress int, speed int64) {
	receiver.Status = status
	receiver.Data = data
	receiver.Progress = progress
	receiver.RunSpeed = speed
	receiver.RunAt = time.Now()
}
