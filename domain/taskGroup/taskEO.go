package taskGroup

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

// TaskEO 任务记录
type TaskEO struct {
	Id             int64                                  // 主键
	TraceId        string                                 // 上下文ID
	Name           string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver            int                                    // 版本
	Caption        string                                 // 任务组标题
	StartAt        dateTime.DateTime                      // 开始时间（计划时间）
	RunAt          dateTime.DateTime                      // 实际执行时间（含结束时间）
	RunSpeed       int64                                  // 运行耗时
	Client         ClientVO                               // 客户端
	Progress       int                                    // 进度0-100
	ScheduleStatus scheduleStatus.Enum                    // 调度状态
	ExecuteStatus  executeStatus.Enum                     // 执行结果
	SchedulerAt    dateTime.DateTime                      // 调度时间
	Data           collections.Dictionary[string, string] // 本次执行任务时的Data数据
	CreateAt       dateTime.DateTime                      // 任务创建时间
	Remark         string                                 // 备注
}

func NewTaskDO() *TaskEO {
	return &TaskEO{}
}

// SetJobName 更新了JobName，则要立即更新Task的JobName
func (receiver *TaskEO) SetJobName(name string) {
	receiver.Name = name
}

// SetFail 设为失败
func (receiver *TaskEO) SetFail(remark string) {
	receiver.ExecuteStatus = executeStatus.Fail
	receiver.Remark = remark
}

// SetScheduling 调度
func (receiver *TaskEO) SetScheduling() {
	receiver.ScheduleStatus = scheduleStatus.Scheduling
}

// ScheduleFail 调度失败
func (receiver *TaskEO) ScheduleFail(remark string) {
	receiver.ScheduleStatus = scheduleStatus.Fail
	receiver.ExecuteStatus = executeStatus.Fail
	receiver.Remark = remark
}

// ScheduleSuccess 调度时设置客户端
func (receiver *TaskEO) ScheduleSuccess(client ClientVO) {
	receiver.ScheduleStatus = scheduleStatus.Success
	receiver.SchedulerAt = dateTime.Now()
	receiver.Client = client
}

// IsNull 未分配
func (receiver *TaskEO) IsNull() bool {
	return receiver.Id == 0 && receiver.Caption == "" && receiver.Name == ""
}

// IsFinish 是否完成
func (receiver *TaskEO) IsFinish() bool {
	return receiver.ExecuteStatus == executeStatus.Success || receiver.ExecuteStatus == executeStatus.Fail
}

// IsWorking 是否为执行中
func (receiver *TaskEO) IsWorking() bool {
	return receiver.ExecuteStatus == executeStatus.Working
}

// UpdateTask 更新任务
func (receiver *TaskEO) UpdateTask(status executeStatus.Enum, data collections.Dictionary[string, string], progress int, remark string) {
	receiver.Data = data
	receiver.Progress = progress
	receiver.UpdateTaskStatus(status, remark)
}

// UpdateTask 更新任务
func (receiver *TaskEO) UpdateTaskStatus(status executeStatus.Enum, remark string) {
	switch status {
	case executeStatus.Fail, executeStatus.Working, executeStatus.Success:
		receiver.ExecuteStatus = status
	default:
		receiver.ExecuteStatus = executeStatus.Fail
		flog.Warningf("任务组 %s %d 回调的状态设置不正确：%d", receiver.Name, receiver.Id, status)
		remark = fmt.Sprintf("回调的状态设置不正确：%d", status)
	}

	receiver.RunAt = dateTime.Now()

	// 耗时
	if status.IsFinish() {
		receiver.RunSpeed = receiver.RunAt.Sub(receiver.SchedulerAt).Milliseconds()
	}

	if remark != "" {
		receiver.Remark = remark
	}

	// 客户端没有设置进度，且执行成功时，自动设为100
	if status == executeStatus.Success {
		receiver.Progress = 100
	}
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		receiver.TraceId = traceContext.GetTraceId()
	}
}
