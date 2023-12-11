// @area /api/
package api

import "FSchedule/domain/client"

// Logout 客户端下线
// @post /logout
func Logout(clientId int64, repository client.Repository) {
	clientDO := repository.ToEntity(clientId)
	clientDO.Logout()
	repository.Save(&clientDO)
}
