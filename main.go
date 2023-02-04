package main

import (
	"FSchedule/application/clientApp"
	"FSchedule/application/taskGroupApp"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func main() {
	fs.Initialize[StartupModule]("FSchedule")
	webapi.Area("/api/", func() {
		// 客户端注册
		webapi.RegisterPOST("/registry", clientApp.Registry)
		// 客户端下线
		webapi.RegisterPOST("/logout", clientApp.Logout)
		// 客户端回调
		webapi.RegisterPOST("/taskReport", taskGroupApp.TaskReport)
		// 上传日志
		webapi.RegisterPOST("/logReport", taskGroupApp.LogReport)
	})
	webapi.UseApiResponse()
	webapi.Run()
}
