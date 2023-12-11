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
    {"GET", "/basicapi/taskGroup/list", basicapi.TaskGroupList, "", []context.IFilter{}, []string{"name", "enable", "taskStatus", "pageSize", "pageIndex", ""}},
    {"GET", "/basicapi/taskGroup/info-{taskGroupId}", basicapi.TaskGroupInfo, "", []context.IFilter{}, []string{"taskGroupId", ""}},
    {"GET", "/basicapi/taskGroup/update", basicapi.TaskGroupUpdate, "", []context.IFilter{}, []string{"req", ""}},
    {"POST", "/basicapi/taskGroup/delete", basicapi.TaskGroupDelete, "", []context.IFilter{}, []string{"taskGroupId", ""}},
    {"GET", "/basicapi/taskGroup/count", basicapi.TaskGroupCount, "", []context.IFilter{}, []string{""}},
    {"GET", "/basicapi/taskGroup/unRunCount", basicapi.TaskGroupUnRunCount, "", []context.IFilter{}, []string{""}},
    {"GET", "/basicapi/taskGroup/schedulerWorkingList", basicapi.TaskGroupSchedulerList, "", []context.IFilter{}, []string{"pageSize", "pageIndex", ""}},
}
