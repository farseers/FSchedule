package enum

type TaskStatus int

const (
	None      TaskStatus = iota // None 未开始
	Scheduler                   // Scheduler 已调度
	Working                     // Working 执行中
	Fail                        // Fail 失败
	Success                     // Success 完成
)

func (e TaskStatus) String() string {
	switch e {
	case None:
		return "None"
	case Scheduler:
		return "Scheduler"
	case Working:
		return "Working"
	case Fail:
		return "Fail"
	case Success:
		return "Success"
	}
	return "None"
}
