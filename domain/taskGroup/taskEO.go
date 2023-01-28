package taskGroup

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
	"time"
)

// TaskEO 任务记录
type TaskEO struct {
	Id          int                                    // 主键
	Ver         int64                                  // 版本
	TaskGroupId int                                    // 任务组ID
	Caption     string                                 // 任务组标题
	JobName     string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	StartAt     time.Time                              // 开始时间
	RunAt       time.Time                              // 实际执行时间
	RunSpeed    int64                                  // 运行耗时
	Client      ClientVO                               // 客户端
	Progress    int                                    // 进度0-100
	Status      enum.TaskStatus                        // 状态
	CreateAt    time.Time                              // 任务创建时间
	SchedulerAt time.Time                              // 调度时间
	Data        collections.Dictionary[string, string] // 本次执行任务时的Data数据
}

func NewTaskDO() *TaskEO {
	return &TaskEO{}
}

// SetClient 调度时设置客户端
func (do *TaskEO) SetClient(client ClientVO) {
	do.Status = enum.Scheduler
	do.SchedulerAt = time.Now()
	do.Client = client
}

// SetJobName 更新了JobName，则要立即更新Task的JobName
func (do *TaskEO) SetJobName(jobName string) {
	do.JobName = jobName
}

// SetFail 设备为失败
func (do *TaskEO) SetFail() {
	do.Status = enum.Fail
}

// IsNull 未分配
func (do *TaskEO) IsNull() bool {
	return do.Id == 0 && do.Status == enum.None && do.Caption == "" && do.TaskGroupId == 0
}

// IsFinish 是否完成
func (do *TaskEO) IsFinish() bool {
	return do.Status == enum.Success || do.Status == enum.Fail
}
