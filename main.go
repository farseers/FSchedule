package main

import (
	"FSchedule/application/clientApp"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func main() {
	fs.Initialize[StartupModule]("FSchedule")
	webapi.Area("/api/", func() {
		webapi.RegisterPOST("/registry", clientApp.Registry)
		webapi.RegisterPOST("/logout", clientApp.Logout)
	})
	webapi.UseApiResponse()
	webapi.Run()
}
