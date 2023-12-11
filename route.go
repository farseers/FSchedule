// 该文件由fsctl route命令自动生成，请不要手动修改此文件
package main

import (
	"FSchedule/application/api"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/context"
)

var route = []webapi.Route{
	{"POST", "/api/logout", api.Logout, "", []context.IFilter{}, []string{"clientId", ""}},
	{"POST", "/api/registry", api.Registry, "", []context.IFilter{}, []string{"dto", "", "", ""}},
	{"POST", "/api/logReport", api.LogReport, "", []context.IFilter{}, []string{"dto", "", ""}},
	{"POST", "/api/taskReport", api.TaskReport, "", []context.IFilter{}, []string{"dto", "", ""}},
}
