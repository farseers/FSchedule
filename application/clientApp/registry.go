// Package clientApp
// @area /api/
package clientApp

import (
	"FSchedule/domain/client"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/webapi"
)

type RegistryDTO struct {
	Id   int64            `json:"ClientId"`   // 客户端ID
	Name string           `json:"ClientName"` // 客户端名称
	Ip   string           `json:"ClientIp"`   // 客户端IP
	Port int              `json:"ClientPort"` // 客户端端口
	Jobs []RegistryJobDTO `json:"ClientJobs"` // 客户端动态注册任务
}

type RegistryJobDTO struct {
	Name     string                                 // 任务名称
	Ver      int                                    // 任务版本
	Caption  string                                 // 任务标题
	Cron     string                                 // 任务执行表达式
	StartAt  int64                                  // 任务开始时间
	IsEnable bool                                   // 任务是否启用
	Data     collections.Dictionary[string, string] // 第一次注册时使用
}

type RegistryResponse struct {
	ClientIp   string // 客户端IP
	ClientPort int    // 客户端端口
}

// Registry 客户端注册
// @post /registry
func Registry(dto RegistryDTO, clientRepository client.Repository, taskGroupRepository taskGroup.Repository, scheduleRepository schedule.Repository) RegistryResponse {
	do := mapper.Single[client.DomainObject](dto)
	// 如果客户端没有指定IP时，由服务端获取
	if do.Ip == "" {
		do.Ip = webapi.GetHttpContext().URI.GetRealIp()
	}
	do.Jobs = collections.NewList[client.JobVO]()
	if do.IsNil() {
		exception.ThrowWebException(403, "客户端ID、Name、IP、Port未完整传入")
	}

	// 更新任务组
	for _, jobDTO := range dto.Jobs {
		lstTaskGroup := taskGroupRepository.ToListByName(jobDTO.Name)
		// 当没有找到任务组时，注册一个新的任务组
		if lstTaskGroup.Count() == 0 {
			taskGroupDO := taskGroup.New(jobDTO.Name, jobDTO.Caption, jobDTO.Ver, jobDTO.Cron, jobDTO.Data, jobDTO.StartAt, jobDTO.IsEnable)
			taskGroupRepository.Save(*taskGroupDO)
		} else { // 找到任务组，则更新现有任务组版本（如果有变动）
			lstTaskGroup.Foreach(func(taskGroupDO *taskGroup.DomainObject) {
				taskGroupDO.UpdateVer(jobDTO.Name, jobDTO.Caption, jobDTO.Ver, jobDTO.Cron, taskGroupDO.Data, jobDTO.StartAt, jobDTO.IsEnable)
				if taskGroupDO.NeedSave {
					taskGroupRepository.Save(*taskGroupDO)
				}
			})
		}
		do.Jobs.Add(mapper.Single[client.JobVO](jobDTO))
	}

	// 保存客户端信息
	do.Registry()
	clientRepository.Save(&do)

	return RegistryResponse{
		ClientIp:   do.Ip,
		ClientPort: do.Port,
	}
}
