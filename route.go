// 该文件由fsctl route命令自动生成，请不要手动修改此文件
package main

import (
	"FSchedule/application/clientApp"
	"FSchedule/application/taskGroupApp"
	"github.com/farseer-go/webapi"
)

var route = []webapi.Route{
	{"POST", "/api/logout", clientApp.Logout, "", nil, []string{"clientId", ""}},
	{"POST", "/api/registry", clientApp.Registry, "", nil, []string{"dto", "", "", ""}},
	{"POST", "/api/logReport", taskGroupApp.LogReport, "", nil, []string{"dto", "", ""}},
	{"POST", "/api/taskReport", taskGroupApp.TaskReport, "", nil, []string{"", "", ""}},
}
