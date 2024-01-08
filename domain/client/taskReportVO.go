package client

import (
	"FSchedule/domain/enum"
	"github.com/farseer-go/collections"
)

type TaskReportVO struct {
	Id           int64                                  // 主键
	Name         string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Data         collections.Dictionary[string, string] // 数据
	NextTimespan int64                                  // 下次执行时间
	Progress     int                                    // 当前进度
	Status       enum.TaskStatus                        // 执行状态
	RunSpeed     int64                                  // 执行速度
}
