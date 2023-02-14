package taskGroup

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
	"time"
)

// TaskEO 任务记录
type TaskEO struct {
	Id          int64                                  // 主键
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
func (do *TaskEO) SetClient(client ClientVO) {
	do.Status = enum.Working
	do.SchedulerAt = time.Now()
	do.Client = client
}

// SetJobName 更新了JobName，则要立即更新Task的JobName
func (do *TaskEO) SetJobName(name string) {
	do.Name = name
}

// SetFail 设备为失败
func (do *TaskEO) SetFail() {
	do.Status = enum.Fail
}

// Scheduling 调度
func (do *TaskEO) Scheduling() {
	do.Status = enum.Scheduling
}

// ScheduleFail 调度失败
func (do *TaskEO) ScheduleFail() {
	do.Status = enum.ScheduleFail
}

// IsNull 未分配
func (do *TaskEO) IsNull() bool {
	return do.Id == 0 && do.Caption == "" && do.Name == ""
}

// IsFinish 是否完成
func (do *TaskEO) IsFinish() bool {
	return do.Status == enum.Success || do.Status == enum.Fail
}

// UpdateTask 更新任务
func (do *TaskEO) UpdateTask(status enum.TaskStatus, data collections.Dictionary[string, string], progress int, speed int64) {
	do.Status = status
	do.Data = data
	do.Progress = progress
	do.RunSpeed = speed
}
