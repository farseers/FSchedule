package client

import (
	"FSchedule/domain/enum/clientStatus"
	"context"
	"time"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/timingWheel"
	"github.com/farseer-go/webapi/websocket"
)

type DomainObject struct {
	Id               string                 // 客户端ID
	Name             string                 // 客户端名称
	Ip               string                 // 客户端IP
	Port             int                    // 客户端端口
	ActivateAt       dateTime.DateTime      // 活动时间
	ScheduleAt       dateTime.DateTime      // 任务调度时间
	Status           clientStatus.Enum      // 客户端状态
	QueueCount       int                    // 排队中的任务数量
	WorkCount        int                    // 正在处理的任务数量
	ErrorCount       int                    // 错误次数
	Job              JobVO                  // 客户端支持的任务
	websocketContext *websocket.BaseContext `json:"-"` // 客户端
	Ctx              context.Context        `json:"-"` // 用于通知应用端是否断开连接
	IsMaster         bool                   // 是否为主客户端
}

// Registry 注册客户端
func (receiver *DomainObject) Registry(websocketContext *websocket.BaseContext) {
	receiver.websocketContext = websocketContext
	receiver.Ctx = websocketContext.Ctx

	receiver.ActivateAt = dateTime.Now()
	receiver.Status = clientStatus.Online

	// 定时保存客户端信息
	receiver.ActivateClient()
}

// IsNil 判断注册的客户端是否有效
func (receiver *DomainObject) IsNil() bool {
	return receiver.Id == "" || receiver.Name == "" || receiver.Ip == "" || receiver.Port == 0
}

// IsOffline 判断客户端是否下线
func (receiver *DomainObject) IsOffline() bool {
	return receiver.Status == clientStatus.Offline
}

// TrySchedule 调度
func (receiver *DomainObject) TrySchedule(task TaskEO) error {
	// 向客户端发起调度请求
	if err := receiver.websocketContext.Send(map[string]any{"Type": 0, "Task": task}); err != nil {
		flog.Warningf("任务组：%s %d 向客户端%s（%s）：%s:%d 调度失败：%s", task.Name, task.Id, receiver.Name, receiver.Id, receiver.Ip, receiver.Port, err.Error())
		receiver.scheduleFail()
		return err
	}

	// 更新客户端状态
	receiver.ScheduleAt = dateTime.Now()
	//receiver.ActivateAt = dateTime.Now()
	receiver.ErrorCount = 0
	receiver.Status = clientStatus.Scheduler
	return nil
}

// 通知客户端，停止任务
func (receiver *DomainObject) Kill(taskId int64) {
	if err := receiver.websocketContext.Send(map[string]any{"Type": 1, "Task": TaskEO{Id: taskId}}); err != nil {
		flog.Warningf("向客户端%s（%s）：%s:%d 停止任务时失败：%s", receiver.Name, receiver.Id, receiver.Ip, receiver.Port, err.Error())
	}
}

// 调度失败
func (receiver *DomainObject) scheduleFail() {
	// 离线状态，不需要设置
	if receiver.IsOffline() {
		return
	}

	// 3次失败，则标记为无法调度
	receiver.ErrorCount++

	// 大于5次、活动时间超过30秒，则判定为离线
	now := dateTime.Now()
	if receiver.ErrorCount >= 5 && now.Sub(receiver.ActivateAt).Seconds() >= 30 {
		receiver.Status = clientStatus.Offline
		flog.Warningf("客户端%s（%s）：%s:%d 调度失败%d次且超过30秒没有活动，状态变更为：%s", receiver.Name, receiver.Id, receiver.Ip, receiver.Port, receiver.ErrorCount, receiver.Status.String())
		receiver.Close()
	}
}

// 关闭客户端
func (receiver *DomainObject) Close() {
	receiver.websocketContext.Close()
}

// 客户端是否关闭
func (receiver *DomainObject) IsClose() bool {
	return receiver.websocketContext == nil || receiver.websocketContext.IsClose()
}

// 心跳检查客户端
func (receiver *DomainObject) ActivateClient() {
	clientRepository := container.Resolve[Repository]()

	clientList.Store(receiver.Id, receiver)
	clientRepository.Save(*receiver)
	flog.Infof("客户端：%s(%s)，%s 连接成功", receiver.Id, receiver.Name, receiver.Job.Name)

	// 定时保存客户端信息
	go func(receiver *DomainObject) {
		defer flog.Infof("客户端：%s(%s)，%s 断开连接", receiver.Id, receiver.Name, receiver.Job.Name)
		for {
			select {
			case <-receiver.Ctx.Done():
				clientList.Delete(receiver.Id)
				clientRepository.RemoveClient(receiver.Id)
				return
			case <-timingWheel.Add(5 * time.Second).C:
				clientDO, _ := clientList.Load(receiver.Id)
				if clientDO == nil || clientDO.(*DomainObject).IsClose() {
					clientList.Delete(receiver.Id)
					clientRepository.RemoveClient(receiver.Id)
					return
				}

				// 发送心跳
				if err := exception.TryCatch(func() {
					err := receiver.websocketContext.Send(map[string]any{"Type": -1})
					exception.ThrowRefuseExceptionError(err)
				}); err != nil {
					//flog.Warningf("向客户端%s（%s）：%s:%d 发送心跳时失败：%s", receiver.Name, receiver.Id, receiver.Ip, receiver.Port, err.Error())
					clientList.Delete(receiver.Id)
					clientRepository.RemoveClient(receiver.Id)
					receiver.websocketContext.Close()
					return
				}

				clientDO.(*DomainObject).ActivateAt = dateTime.Now()
				clientRepository.Save(*clientDO.(*DomainObject))
			}
		}
	}(receiver)
}
