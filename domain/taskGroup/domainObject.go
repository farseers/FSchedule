package taskGroup

import (
	"github.com/farseer-go/collections"
	"time"
)

type DomainObject struct {
	Name        string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver         int                                    // 版本
	Task        TaskEO                                 // 任务
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
func (receiver *DomainObject) UpdateVer(name string, caption string, ver int, cron string, StartAt int64, enable bool) {
	// 只更新高一个版本号的数据
	if receiver.Ver+1 == ver {
		receiver.Name = name
		receiver.Caption = caption
		receiver.Ver = ver
		receiver.Cron = cron
		receiver.StartAt = time.Unix(StartAt, 0)
		receiver.NeedSave = true
		receiver.IsEnable = enable
	}
}
