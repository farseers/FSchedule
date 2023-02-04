package clientApp

import (
	"FSchedule/domain/client"
	"FSchedule/domain/taskGroup"
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
	Caption  string // 任务标题
	Ver      int    // 任务版本
	Cron     string // 任务执行表达式
	StartAt  int64  // 任务开始时间
	IsEnable bool   // 任务开始时间
}

// Registry 客户端注册
func Registry(dto RegistryDTO, clientRepository client.Repository, taskGroupRepository taskGroup.Repository) {
	do := mapper.Single[client.DomainObject](dto)
	if do.IsNil() {
		exception.ThrowWebException(403, "客户端ID、Name、IP、Port未完整传入")
	}

	// 保存客户端信息
	do.Registry()
	clientRepository.Save(do)

	// 更新任务组
	for _, jobDTO := range dto.Jobs {
		taskGroupDO := taskGroupRepository.ToEntity(jobDTO.Name)
		taskGroupDO.UpdateVer(jobDTO.Name, jobDTO.Caption, jobDTO.Ver, jobDTO.Cron, jobDTO.StartAt, jobDTO.IsEnable)
		if taskGroupDO.NeedSave {
			taskGroupRepository.Save(taskGroupDO)
		}
	}
}
