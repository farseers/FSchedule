package clientStatus

type Enum int

const (
	Online       Enum = iota // 刚上线
	Scheduler                // 接受调度
	UnSchedule               // 无法调度（请求出错）
	StopSchedule             // 拒绝调度（客户端在忙）
	Offline                  // 离线
)

func (e Enum) String() string {
	switch e {
	case Online:
		return "刚上线"
	case Scheduler:
		return "接受调度"
	case UnSchedule:
		return "无法调度"
	case StopSchedule:
		return "拒绝调度"
	case Offline:
		return "离线"
	}
	return "刚上线"
}
