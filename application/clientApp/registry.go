package clientApp

import (
	"FSchedule/domain/client"
	"FSchedule/domain/schedule"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

type RegistryDTO struct {
	Id   int64            `json:"ClientId"`   // 客户端ID
	Name string           `json:"ClientName"` // 客户端名称
	Ip   string           `json:"ClientIp"`   // 客户端IP
	Port int              `json:"ClientPort"` // 客户端端口
	Jobs []RegistryJobDTO `json:"ClientJobs"` // 客户端动态注册任务
}

type RegistryJobDTO struct {
	Name     string // 任务名称
	Ver      int    // 任务版本
	Caption  string // 任务标题
	Cron     string // 任务执行表达式
	StartAt  int64  // 任务开始时间
	IsEnable bool   // 任务是否启用
}

// Registry 客户端注册
func Registry(dto RegistryDTO, clientRepository client.Repository, taskGroupRepository taskGroup.Repository, scheduleRepository schedule.Repository) {
	do := mapper.Single[client.DomainObject](dto)
	do.Jobs = collections.NewList[client.JobVO]()
	if do.IsNil() {
		exception.ThrowWebException(403, "客户端ID、Name、IP、Port未完整传入")
	}

	// 先推送任务信息再保存客户端
	// 更新任务组
	for _, jobDTO := range dto.Jobs {
		taskGroupDO := taskGroupRepository.ToEntity(jobDTO.Name)
		taskGroupDO.UpdateVer(jobDTO.Name, jobDTO.Caption, jobDTO.Ver, jobDTO.Cron, jobDTO.StartAt, jobDTO.IsEnable)
		if taskGroupDO.NeedSave {
			taskGroupRepository.Save(taskGroupDO)
		}
		do.Jobs.Add(mapper.Single[client.JobVO](jobDTO))
	}

	// 保存客户端信息
	do.Registry()
	do.CheckOnline()
	clientRepository.Save(&do)
}
