package http

import (
	"FSchedule/application/taskGroupApp"
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
		return client.ResourceVO{}, err
	}
	if apiResponse.StatusCode != 200 {
		log := fmt.Sprintf("客户端：http://%s:%d，状态码：%d，错误内容：%s", do.Ip, do.Port, apiResponse.StatusCode, apiResponse.StatusMessage)
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
		log := fmt.Sprintf("客户端：http://%s:%d，状态码：%d，错误内容：%s", do.Ip, do.Port, apiResponse.StatusCode, apiResponse.StatusMessage)
		flog.Info(log)
		return client.ResourceVO{}, flog.Error(log)
	}
	return apiResponse.Data, nil
}

func (receiver clientHttp) Status(do *client.DomainObject, taskId int64) (taskGroupApp.TaskReportDTO, error) {
	clientUrl := fmt.Sprintf("http://%s:%d/api/status", do.Ip, do.Port)
	var apiResponse core.ApiResponse[taskGroupApp.TaskReportDTO]
	body := map[string]any{
		"taskId": taskId,
	}
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(body).PostUnmarshal(&apiResponse)
	if err != nil {
		return taskGroupApp.TaskReportDTO{}, err
	}
	if apiResponse.StatusCode != 200 {
		log := fmt.Sprintf("客户端：http://%s:%d，状态码：%d，错误内容：%s", do.Ip, do.Port, apiResponse.StatusCode, apiResponse.StatusMessage)
		flog.Info(log)
		return taskGroupApp.TaskReportDTO{}, flog.Error(log)
	}
	return apiResponse.Data, nil
}

func (receiver clientHttp) Kill(do *client.DomainObject, taskId int64) bool {
	clientUrl := fmt.Sprintf("http://%s:%d/api/kill", do.Ip, do.Port)
	var apiResponse core.ApiResponse[any]
	body := map[string]any{
		"taskId": taskId,
	}
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(body).PostUnmarshal(&apiResponse)
	if err != nil {
		return false
	}
	if apiResponse.StatusCode != 200 {
		flog.Infof("客户端：http://%s:%d，状态码：%s", do.Ip, do.Port, apiResponse.StatusCode)
		return false
	}
	return true
}
