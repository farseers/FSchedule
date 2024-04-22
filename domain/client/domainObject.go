package client

import (
	"FSchedule/domain/enum/clientStatus"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
)

type DomainObject struct {
	Id          int64                   // 客户端ID
	Name        string                  // 客户端名称
	Ip          string                  // 客户端IP
	Port        int                     // 客户端端口
	ActivateAt  dateTime.DateTime       // 活动时间
	ScheduleAt  dateTime.DateTime       // 任务调度时间
	Status      clientStatus.Enum       // 客户端状态
	QueueCount  int                     // 排队中的任务数量
	WorkCount   int                     // 正在处理的任务数量
	CpuUsage    float64                 // CPU百分比
	MemoryUsage float64                 // 内存百分比
	ErrorCount  int                     // 错误次数
	Jobs        collections.List[JobVO] // 客户端支持的任务
	NeedNotice  bool                    //	是否需要通知任务组
}

// IsNil 判断注册的客户端是否有效
func (receiver *DomainObject) IsNil() bool {
	return receiver.Id == 0 || receiver.Name == "" || receiver.Ip == "" || receiver.Port == 0
}

// IsOnline 是否刚注册进来
func (receiver *DomainObject) IsOnline() bool {
	return receiver.Status == clientStatus.Online
}

// IsOffline 判断客户端是否下线
func (receiver *DomainObject) IsOffline() bool {
	return receiver.Status == clientStatus.Offline
}

// IsNotSchedule 状态不是调度状态
func (receiver *DomainObject) IsNotSchedule() bool {
	return receiver.Status != clientStatus.Scheduler
}

// IsNotSchedule 可调度状态
func (receiver *DomainObject) IsCanSchedule() bool {
	return receiver.Status == clientStatus.Scheduler || receiver.Status == clientStatus.Online
}

// Registry 注册客户端
func (receiver *DomainObject) Registry() {
	receiver.ActivateAt = dateTime.Now()
	receiver.Status = clientStatus.Online
	receiver.NeedNotice = true
}

// Logout 客户端下线
func (receiver *DomainObject) Logout() {
	receiver.Status = clientStatus.Offline
	receiver.NeedNotice = true
}

// CheckOnline 检查客户端是否存活
func (receiver *DomainObject) CheckOnline() error {
	status, err := container.Resolve[IClientCheck]().Check(receiver)
	if err != nil {
		flog.Warningf("检查客户端%s（%d）：%s:%d 在线状态失败：%s", receiver.Name, receiver.Id, receiver.Ip, receiver.Port, err.Error())
		receiver.scheduleFail()
		return err
	}
	// 检查成功
	receiver.updateStatus(status)
	return nil
}

// 向客户端检查任务状态
func (receiver *DomainObject) CheckTaskStatus(taskGroupName string, taskId int64) (TaskReportVO, error) {
	clientCheck := container.Resolve[IClientCheck]()

	dto, err := clientCheck.Status(receiver, taskGroupName, taskId)
	if err != nil {
		flog.Warningf("向客户端%s（%d）：%s:%d 检查任务失败：%s", receiver.Name, receiver.Id, receiver.Ip, receiver.Port, err.Error())
		receiver.scheduleFail()
		return TaskReportVO{}, err
	} else {
		receiver.updateStatus(dto.ResourceVO)
	}
	return dto, nil
}

// TrySchedule 调度
func (receiver *DomainObject) TrySchedule(task TaskEO) (bool, error) {
	var status ResourceVO
	var err error

	// 向客户端发起调度请求
	if status, err = container.Resolve[IClientCheck]().Invoke(receiver, task); err != nil {
		flog.Warningf("任务组：%s %d 向客户端%s（%d）：%s:%d 调度失败：%s", task.Name, task.Id, receiver.Name, receiver.Id, receiver.Ip, receiver.Port, err.Error())
		receiver.scheduleFail()
		return false, err
	}

	// 调度成功
	receiver.ScheduleAt = dateTime.Now()
	receiver.updateStatus(status)
	return true, nil
}

// 更新客户端状态
func (receiver *DomainObject) updateStatus(status ResourceVO) {
	oldStatus := receiver.Status
	receiver.ActivateAt = dateTime.Now()
	receiver.ErrorCount = 0
	receiver.CpuUsage = status.CpuUsage
	receiver.MemoryUsage = status.MemoryUsage
	receiver.QueueCount = status.QueueCount
	receiver.WorkCount = status.WorkCount

	if status.AllowSchedule {
		receiver.Status = clientStatus.Scheduler
	} else {
		receiver.Status = clientStatus.StopSchedule
	}

	receiver.NeedNotice = oldStatus != receiver.Status
}

// 调度失败
func (receiver *DomainObject) scheduleFail() {
	// 离线状态，不需要设置
	if receiver.IsOffline() {
		return
	}

	// 3次失败，则标记为无法调度
	receiver.ErrorCount++
	if receiver.ErrorCount >= 3 {
		receiver.Status = clientStatus.UnSchedule
		flog.Warningf("客户端%s（%d）：%s:%d 调度失败%d次，状态变更为：%s", receiver.Name, receiver.Id, receiver.Ip, receiver.Port, receiver.ErrorCount, receiver.Status.String())
	}

	// 大于5次、活动时间超过30秒，则判定为离线
	now := dateTime.Now()
	if receiver.ErrorCount >= 5 && now.Sub(receiver.ActivateAt).Seconds() >= 30 {
		receiver.Logout()
		flog.Warningf("客户端%s（%d）：%s:%d 调度失败%d次且超过30秒没有活动，状态变更为：%s", receiver.Name, receiver.Id, receiver.Ip, receiver.Port, receiver.ErrorCount, receiver.Status.String())
	}
}
