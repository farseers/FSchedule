package clientApp

import (
	"FSchedule/domain/client"
	"FSchedule/domain/taskGroup"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

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
