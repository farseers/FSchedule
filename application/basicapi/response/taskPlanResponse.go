package response

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
)

type TaskPlanResponse struct {
	Id             string                                 // 主键
	TraceId        string                                 // 上下文ID
	Name           string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver            int                                    // 版本
	Caption        string                                 // 任务组标题
	Plan           string                                 // 开始时间
	RunAt          dateTime.DateTime                      // 实际执行时间
	RunSpeed       int64                                  // 运行耗时
	Client         taskGroup.ClientVO                     // 客户端
	Progress       int                                    // 进度0-100
	ScheduleStatus scheduleStatus.Enum                    // 调度状态
	ExecuteStatus  executeStatus.Enum                     // 状态
	SchedulerAt    dateTime.DateTime                      // 调度时间
	Data           collections.Dictionary[string, string] // 本次执行任务时的Data数据
	CreateAt       dateTime.DateTime                      // 任务创建时间
}
