package main

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/webapi"
)

func main() {
	fs.Initialize[StartupModule]("FSchedule")
	webapi.RegisterRoutes(route...)
	webapi.UseApiResponse()
	webapi.UsePprof()
	webapi.UseApiDoc()
	webapi.UseValidate()
	webapi.UseCors()
	webapi.PrintRoute()
	webapi.Run()
}
