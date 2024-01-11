package enum

type ClientStatus int

const (
	Online       ClientStatus = iota // 刚上线
	Scheduler                        // 接受调度
	UnSchedule                       // 无法调度（请求出错）
	StopSchedule                     // 拒绝调度（客户端在忙）
	Offline                          // 离线
)
