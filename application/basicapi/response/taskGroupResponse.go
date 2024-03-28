package response

import (
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
)

type TaskGroupResponse struct {
	Name        string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Ver         int                                    // 版本
	Task        taskGroup.TaskEO                       // 最新的任务
	Caption     string                                 // 任务组标题
	Data        collections.Dictionary[string, string] // 本次执行任务时的Data数据
	StartAt     dateTime.DateTime                      // 开始时间
	NextAt      dateTime.DateTime                      // 下次执行时间
	Cron        string                                 // 时间定时器表达式
	ActivateAt  dateTime.DateTime                      // 活动时间
	LastRunAt   dateTime.DateTime                      // 最后一次完成时间
	IsEnable    bool                                   // 是否开启
	RunSpeedAvg int64                                  // 运行平均耗时
	RunCount    int                                    // 运行次数
	Clients     collections.List[ClientResponse]       // 客户端列表
}

type ClientResponse struct {
	Id   string // 客户端ID
	Name string // 客户端名称
	Ip   string // 客户端IP
	Port int    // 客户端端口
}
