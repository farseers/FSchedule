package client

// ResourceVO 客户端资源情况
type ResourceVO struct {
	QueueCount    int     // 排队中的任务数量
	WorkCount     int     // 正在处理的任务数量
	CpuUsage      float32 // CPU百分比
	MemoryUsage   float32 // 内存百分比
	AllowSchedule bool    // 是否允许调度
}
