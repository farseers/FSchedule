// 该文件由fsctl route命令自动生成，请不要手动修改此文件
package main

import (
	"FSchedule/application/clientApp"
	"FSchedule/application/taskGroupApp"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/context"
)

var route = []webapi.Route{
	{"POST", "/api/logout", clientApp.Logout, "", []context.IFilter{}, []string{"clientId", ""}},
	{"POST", "/api/registry", clientApp.Registry, "", []context.IFilter{}, []string{"dto", "", "", ""}},
	{"POST", "/api/logReport", taskGroupApp.LogReport, "", []context.IFilter{}, []string{"dto", "", ""}},
	{"POST", "/api/taskReport", taskGroupApp.TaskReport, "", []context.IFilter{}, []string{"dto", "", ""}},
}
