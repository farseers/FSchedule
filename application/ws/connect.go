// @area /ws/
package ws

import (
	"FSchedule/domain"
	"FSchedule/domain/client"
	"FSchedule/domain/taskGroup"
	"FSchedule/domain/taskLog"
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/webapi/websocket"
	"strings"
)

type request struct {
	Type       int // 接收消息的类型
	Registry   domain.RegistryDTO
	TaskReport client.TaskReportVO
	Log        struct {
		TaskId   int64                                  // 主键
		Ver      int                                    // 版本
		Name     string                                 // 实现Job的特性名称（客户端识别哪个实现类）
		Caption  string                                 // 任务标题
		Data     collections.Dictionary[string, string] // 本次执行任务时的Data数据
		LogLevel eumLogLevel.Enum
		CreateAt int64
		Content  string
	}
}

const tokenName = "FSS-ACCESS-TOKEN"

// 客户端请求任务组分派任务，客户端每个任务单独连接
// @ws connect
func Connect(wsContext *websocket.Context[request], clientRepository client.Repository, taskGroupRepository taskGroup.Repository, taskLogRepository taskLog.Repository) {
	// 新的客户端、先做注册
	req := wsContext.Receiver()
	addr := strings.Split(wsContext.HttpContext.URI.RemoteAddr, ":")
	req.Registry.ClientIp = addr[0]
	req.Registry.ClientPort = parse.ToInt(addr[1])
	clientId := fmt.Sprintf("%s:%d", req.Registry.ClientIp, req.Registry.ClientPort)

	// 客户端注册
	domain.Registry(wsContext.BaseContext, clientId, req.Registry, clientRepository, taskGroupRepository)

	for {
		req = wsContext.Receiver()
		exception.Try(func() {
			switch req.Type {
			case 0: // 客户端回调
				domain.TaskReportService(clientId, req.TaskReport, taskGroupRepository)
			case 1: // 日志上报
				taskLogDO := taskLog.NewDO(req.Log.Name, req.Log.Caption, req.Log.Ver, req.Log.TaskId, req.Log.Data, req.Log.LogLevel, req.Log.Content, req.Log.CreateAt)
				taskLogRepository.Add(taskLogDO)
			}
		}).CatchException(func(exp any) {
			_ = wsContext.Send(core.ApiResponseIntError(fmt.Sprint(exp), 500))
		})
	}
}
