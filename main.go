package main

import (
	"FSchedule/application/clientApp"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func main() {
	fs.Initialize[StartupModule]("FSchedule")
	webapi.RegisterPOST("/registry", clientApp.Registry)
	webapi.UseApiResponse()
	webapi.Run()
}
