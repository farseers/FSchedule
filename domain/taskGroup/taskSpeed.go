package taskGroup

import (
	"github.com/farseer-go/collections"
)

// TaskSpeed 任务执行速度
type TaskSpeed struct {
	speedList []int64 // 任务的所有执行速度
}

func NewTaskSpeed(speedList []int64) *TaskSpeed {
	return &TaskSpeed{speedList: speedList}
}

// GetAvgSpeed 任务的执行平均速度
func (receiver *TaskSpeed) GetAvgSpeed() int64 {
	if len(receiver.speedList) == 0 {
		return 0
	}
	sum := collections.NewList(receiver.speedList...).SumItem()
	return sum / int64(len(receiver.speedList))
}
