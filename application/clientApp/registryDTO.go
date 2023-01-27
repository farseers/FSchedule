package clientApp

type RegistryDTO struct {
	Id   int64            `json:"ClientId"`   // 客户端ID
	Name string           `json:"ClientName"` // 客户端名称
	Ip   string           `json:"ClientIp"`   // 客户端IP
	Port int              `json:"ClientPort"` // 客户端端口
	Jobs []RegistryJobDTO `json:"ClientJobs"` // 客户端动态注册任务
}

type RegistryJobDTO struct {
	Name    string // 任务名称
	Caption string // 任务标题
	Ver     int    // 任务版本
	Cron    string // 任务执行表达式
	StartAt int64  // 任务开始时间
}