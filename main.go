package main

import (
	"fss/domain/clients/client"
	"fss/interfaces"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/webapi"
)

func main() {
	fs.Initialize[StartupModule]("fss")

	webapi.RegisterPOST("test/", func() collections.List[client.DomainObject] {
		repository := container.Resolve[client.Repository]()
		return repository.ToList()
	})

	webapi.RegisterRoutes(routeMeta)
	webapi.RegisterController(&interfaces.TaskController{})

	webapi.UseApiResponse()
	webapi.Run()
}
