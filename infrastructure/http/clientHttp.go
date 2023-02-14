package http

import (
	"FSchedule/domain/client"
	"fmt"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/http"
)

const tokenName = "FSS-ACCESS-TOKEN"

var token = configure.GetString("FSchedule.Server.Token")

type clientHttp struct {
}

func (receiver clientHttp) Check(do *client.DomainObject) (client.ResourceVO, error) {
	clientUrl := fmt.Sprintf("http://%s:%d/api/check", do.Ip, do.Port)
	var apiResponse core.ApiResponse[client.ResourceVO]
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).PostUnmarshal(&apiResponse)
	if err != nil {
		flog.Warningf("客户端（%d）：%s:%d  检查失败", do.Id, do.Ip, do.Port)
		return client.ResourceVO{}, err
	}
	if apiResponse.StatusCode != 200 {
		log := fmt.Sprintf("客户端（%d）：%s，状态码：%d，错误内容：%s", do.Id, clientUrl, apiResponse.StatusCode, apiResponse.StatusMessage)
		flog.Info(log)
		return client.ResourceVO{}, flog.Error(log)
	}
	return apiResponse.Data, nil
}

func (receiver clientHttp) Invoke(do *client.DomainObject, task *client.TaskEO) (client.ResourceVO, error) {
	clientUrl := fmt.Sprintf("http://%s:%d/api/invoke", do.Ip, do.Port)
	var apiResponse core.ApiResponse[client.ResourceVO]
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(task).PostUnmarshal(&apiResponse)
	if err != nil {
		return client.ResourceVO{}, err
	}
	if apiResponse.StatusCode != 200 {
		log := fmt.Sprintf("客户端（%d）：%s，状态码：%d，错误内容：%s", do.Id, clientUrl, apiResponse.StatusCode, apiResponse.StatusMessage)
		flog.Info(log)
		return client.ResourceVO{}, flog.Error(log)
	}
	return apiResponse.Data, nil
}

func (receiver clientHttp) Status(do *client.DomainObject, taskId int64) (client.TaskReportVO, error) {
	clientUrl := fmt.Sprintf("http://%s:%d/api/status", do.Ip, do.Port)
	var apiResponse core.ApiResponse[client.TaskReportVO]
	body := map[string]any{
		"TaskId": taskId,
	}
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(body).PostUnmarshal(&apiResponse)
	if err != nil {
		return client.TaskReportVO{}, err
	}
	if apiResponse.StatusCode != 200 {
		log := fmt.Sprintf("客户端（%d）：%s，状态码：%d，错误内容：%s", do.Id, clientUrl, apiResponse.StatusCode, apiResponse.StatusMessage)
		flog.Info(log)
		return client.TaskReportVO{}, flog.Error(log)
	}
	return apiResponse.Data, nil
}

func (receiver clientHttp) Kill(do *client.DomainObject, taskId int64) bool {
	clientUrl := fmt.Sprintf("http://%s:%d/api/kill", do.Ip, do.Port)
	var apiResponse core.ApiResponse[any]
	body := map[string]any{
		"TaskId": taskId,
	}
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(body).PostUnmarshal(&apiResponse)
	if err != nil {
		return false
	}
	if apiResponse.StatusCode != 200 {
		log := fmt.Sprintf("客户端（%d）：%s，状态码：%d，错误内容：%s", do.Id, clientUrl, apiResponse.StatusCode, apiResponse.StatusMessage)
		flog.Info(log)
		return false
	}
	return true
}
