package domain

import (
	"FSchedule/domain/client"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/webapi/websocket"
)

type RegistryDTO struct {
	ClientName string   // 客户端名称
	ClientIp   string   // 客户端IP
	ClientPort int      // 客户端端口
	Job        struct { // 客户端动态注册任务
		Name     string                                 // 任务名称
		Ver      int                                    // 任务版本
		Caption  string                                 // 任务标题
		Cron     string                                 // 任务执行表达式
		StartAt  int64                                  // 任务开始时间
		IsEnable bool                                   // 任务是否启用
		Data     collections.Dictionary[string, string] // 第一次注册时使用
	}
}

// Registry 客户端注册
func Registry(websocketContext *websocket.BaseContext, clientId string, dto RegistryDTO, clientRepository client.Repository, taskGroupRepository taskGroup.Repository) {
	if dto.ClientName == "" || dto.Job.Name == "" {
		exception.ThrowWebExceptionf(403, "ClientName=%s、JobName=%s，未完整传入", dto.ClientName, dto.Job.Name)
	}

	// 确认cron格式是否正确
	_, err := taskGroup.StandardParser.Parse(dto.Job.Cron)
	if err != nil {
		exception.ThrowWebExceptionf(403, "任务组:%s %s，Cron格式[%s]错误:%s", dto.Job.Name, dto.Job.Caption, dto.Job.Cron, err.Error())
	}

	// 新增 或 修改任务组
	taskGroupDO := taskGroupRepository.ToEntity(dto.Job.Name)
	// 当没有找到任务组时，注册一个新的任务组
	if taskGroupDO.IsNil() {
		taskGroupDO = taskGroup.New(dto.Job.Name, dto.Job.Caption, dto.Job.Ver, dto.Job.Cron, dto.Job.Data, dto.Job.StartAt, dto.Job.IsEnable)
		taskGroupRepository.Save(taskGroupDO)
	} else {
		// 找到任务组，则更新现有任务组版本（如果有变动）
		taskGroupDO.UpdateVer(dto.Job.Name, dto.Job.Caption, dto.Job.Ver, dto.Job.Cron, dto.Job.Data, dto.Job.StartAt, dto.Job.IsEnable)
		if taskGroupDO.NeedSave {
			taskGroupRepository.Save(taskGroupDO)
		}
	}

	clientDO := &client.DomainObject{
		Id:   clientId,
		Ip:   dto.ClientIp,
		Port: dto.ClientPort,
		Name: dto.ClientName,
		Job:  mapper.Single[client.JobVO](dto.Job),
	}

	clientDO.Registry(websocketContext)

	// 将任务组接入监控
	MonitorTaskGroupPush(clientDO, &taskGroupDO)
}
