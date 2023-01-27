package clientApp

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

// Registry 客户端注册
func Registry(dto RegistryDTO, repository client.Repository) {
	do := mapper.Single[client.DomainObject](dto)
	if do.IsNil() {
		exception.ThrowWebException(403, "客户端ID、Name、IP、Port未完整传入")
	}

	do.Registry()

	// 保存客户端信息
	repository.Save(do)
}

// Logout 客户端下线
func Logout(clientId int64, repository client.Repository) {
	repository.ToEntity(clientId).Logout()
}
