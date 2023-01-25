package enum

type EumClientStatus int

const (
	Online       EumClientStatus = iota // 刚上线
	Scheduling                          // 接受调度
	StopSchedule                        // 拒绝调度
	Offline                             // 离线
)
