package request

import (
	"github.com/farseer-go/collections"
	"time"
)

type TaskGroupUpdateRequest struct {
	Id       int64                                  // 主键ID
	Name     string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver      int                                    // 版本
	Caption  string                                 // 任务组标题
	Data     collections.Dictionary[string, string] // 本次执行任务时的Data数据
	StartAt  time.Time                              // 开始时间
	NextAt   time.Time                              // 下次执行时间
	Cron     string                                 // 时间定时器表达式
	IsEnable bool                                   // 是否开启
}