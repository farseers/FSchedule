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
		return "None"
	case Scheduling:
		return "Scheduling"
	case Success:
		return "Success"
	case Fail:
		return "Fail"
	}
	return "None"
}
