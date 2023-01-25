package clientApp

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"time"
)

// Registry 客户端注册
func Registry(dto RegistryDTO, repository client.Repository) {
	do := mapper.Single[client.DomainObject](dto)
	if do.IsNil() {
		exception.ThrowWebException(403, "客户端ID、Name、IP、Port未完整传入")
	}

	do.ActivateAt = time.Now()
	// 保存客户端信息
	repository.Save(do)
}
