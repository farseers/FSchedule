package scheduleStatus

type Enum int

const (
	None       Enum = iota //  未调度
	Scheduling             //  调度中（即将请求客户端）
	Success                //  调度成功
	Fail                   //  调度失败
)

func (e Enum) String() string {
	switch e {
	case None:
		return "未调度"
	case Scheduling:
		return "调度中"
	case Success:
		return "调度成功"
	case Fail:
		return "调度失败"
	}
	return "未调度"
}
