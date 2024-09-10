package client

// ResourceVO 客户端资源情况
type ResourceVO struct {
	QueueCount    int  // 排队中的任务数量
	WorkCount     int  // 正在处理的任务数量
	AllowSchedule bool // 是否允许调度
}
