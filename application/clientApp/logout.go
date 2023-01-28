package clientApp

import "FSchedule/domain/client"

// Logout 客户端下线
func Logout(clientId int64, repository client.Repository) {
	repository.ToEntity(clientId).Logout()
}
