// 该文件由fsctl route命令自动生成，请不要手动修改此文件
package main

import (
	"FSchedule/application/api"
	"FSchedule/application/basicapi"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/context"
)

var route = []webapi.Route{
	{"POST", "/api/logReport", api.LogReport, "", []context.IFilter{}, []string{"dto", "", ""}},
	{"POST", "/api/logout", api.Logout, "", []context.IFilter{}, []string{"clientId", ""}},
	{"POST", "/api/registry", api.Registry, "", []context.IFilter{}, []string{"dto", "", "", ""}},
	{"POST", "/api/taskReport", api.TaskReport, "", []context.IFilter{}, []string{"dto", "", ""}},
	{"POST", "/api/killTask", api.KillTask, "", []context.IFilter{}, []string{"taskGroupName", "", "", ""}},
	{"GET", "/basicapi/client/list", basicapi.ClientList, "", []context.IFilter{}, []string{""}},
	{"GET", "/basicapi/log/list", basicapi.LogList, "", []context.IFilter{}, []string{"taskGroupName", "logLevel", "taskId", "pageSize", "pageIndex", ""}},
	{"GET", "/basicapi/log/listByClientName", basicapi.LogListByClientName, "", []context.IFilter{}, []string{"clientName", "taskGroupName", "logLevel", "taskId", "pageSize", "pageIndex", ""}},
	{"GET", "/basicapi/server/list", basicapi.ServerList, "", []context.IFilter{}, []string{""}},
	{"GET", "/basicapi/stat/statList", basicapi.StatList, "", []context.IFilter{}, []string{""}},
	{"GET", "/basicapi/stat/info", basicapi.Info, "", []context.IFilter{}, []string{""}},
	{"GET", "/basicapi/task/list", basicapi.TaskList, "", []context.IFilter{}, []string{"clientName", "taskGroupName", "scheduleStatus", "executeStatus", "taskId", "pageSize", "pageIndex", ""}},
	{"GET", "/basicapi/task/planList", basicapi.TaskPlanList, "", []context.IFilter{}, []string{"top", ""}},
	{"GET", "/basicapi/taskGroup/list", basicapi.TaskGroupList, "", []context.IFilter{}, []string{"clientName", "taskGroupName", "enable", "taskStatus", "taskId", "clientId", "pageSize", "pageIndex", ""}},
	{"GET", "/basicapi/taskGroup/info-{taskGroupName}", basicapi.TaskGroupInfo, "", []context.IFilter{}, []string{"taskGroupName", ""}},
	{"POST", "/basicapi/taskGroup/update", basicapi.TaskGroupUpdate, "", []context.IFilter{}, []string{"req", ""}},
	{"POST", "/basicapi/taskGroup/delete", basicapi.TaskGroupDelete, "", []context.IFilter{}, []string{"taskGroupName", ""}},
	{"POST", "/basicapi/taskGroup/setEnable", basicapi.SetEnable, "", []context.IFilter{}, []string{"taskGroupName", "enable", ""}},
}
