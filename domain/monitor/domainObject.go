package monitor

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/fs/core"
	"time"
)

// DomainObject 等待任务执行
type DomainObject struct {
	Name     string          // 实现Job的特性名称（客户端识别哪个实现类）
	Ver      int             // 版本
	StartAt  time.Time       // 开始时间
	NextAt   time.Time       // 下次执行时间
	Cron     string          // 时间定时器表达式
	IsEnable bool            // 是否开启
	Status   enum.TaskStatus // 状态
	EventBus core.IEvent     `inject:"TaskSchedule"` // 任务调度事件
}

// Start 监听任务组
func (receiver *DomainObject) Start(c chan DomainObject) {
	for {
		// 任务还没有达到开始之前的准备阶段：
		// 1、开启状态。
		// 2、开始时间 < 当前时间。
		// 3、任务状态=None
		receiver.waitStart(c)

		// 等待时间达了之后，开始调度
		receiver.waitNextAt(c)
	}
}

// 等待开始
func (receiver *DomainObject) waitStart(c chan DomainObject) {
	for {
		select {
		case <-time.After(receiver.StartAt.Sub(time.Now())): // 时间到了，可以开始计算任务执行赶时间
			// 开启状态，且未调度
			if receiver.IsEnable && receiver.Status == enum.None {
				return
			}
			// 等待更新
			receiver.update(<-c)
		case newData := <-c: // 有更新
			receiver.update(newData)
		}
	}
}

// 等待任务执行时间
func (receiver *DomainObject) waitNextAt(c chan DomainObject) {
	select {
	case <-time.After(receiver.NextAt.Sub(time.Now())): // 时间到了，需要调度
		// 标记为调度中，阻止当前监听逻辑重复执行，否则会不停的重复执行调度
		receiver.Status = enum.Scheduler
		receiver.EventBus.Publish(receiver)
	case newData := <-c: // 有更新
		receiver.update(newData)
	}
}

// 有更新
func (receiver *DomainObject) update(newData DomainObject) {
	receiver.Ver = newData.Ver
	receiver.StartAt = newData.StartAt
	receiver.NextAt = newData.NextAt
	receiver.Cron = newData.Cron
	receiver.IsEnable = newData.IsEnable
}
