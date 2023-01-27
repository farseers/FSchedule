package http

import (
	"FSchedule/domain/client"
	"github.com/farseer-go/fs/container"
)

// InitHttp 初始化客户端
func InitHttp() {
	// 注册仓储
	container.Register(func() client.IClientCheck {
		return &clientHttp{}
	})
}
