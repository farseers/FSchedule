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
	body := map[string]any{
		"clientId": do.Id,
	}
	var apiResponse core.ApiResponse[client.ResourceVO]
	_, err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(body).PostUnmarshal(&apiResponse)
	if err != nil {
		return client.ResourceVO{}, err
	}
	if apiResponse.StatusCode != 200 {
		log := fmt.Sprintf("状态码：%d，错误内容：%s", apiResponse.StatusCode, apiResponse.StatusMessage)
		return client.ResourceVO{}, flog.Error(log)
	}
	return apiResponse.Data, nil
}

func (receiver clientHttp) Invoke(do *client.DomainObject, task client.TaskEO) (client.ResourceVO, error) {
	clientUrl := fmt.Sprintf("http://%s:%d/api/invoke", do.Ip, do.Port)
	var apiResponse core.ApiResponse[client.ResourceVO]
	_, err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(task).PostUnmarshal(&apiResponse)
	if err != nil {
		return client.ResourceVO{}, err
	}
	if apiResponse.StatusCode != 200 {
		log := fmt.Sprintf("状态码：%d，错误内容：%s", apiResponse.StatusCode, apiResponse.StatusMessage)
		return client.ResourceVO{}, flog.Error(log)
	}
	return apiResponse.Data, nil
}

func (receiver clientHttp) Status(do *client.DomainObject, taskGroupName string, taskId int64) (client.TaskReportVO, error) {
	clientUrl := fmt.Sprintf("http://%s:%d/api/status", do.Ip, do.Port)
	var apiResponse core.ApiResponse[client.TaskReportVO]
	body := map[string]any{
		"TaskId":     taskId,
		"ClientName": taskGroupName,
	}
	_, err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(body).PostUnmarshal(&apiResponse)
	if err != nil {
		return client.TaskReportVO{}, err
	}
	if apiResponse.StatusCode != 200 {
		log := fmt.Sprintf("状态码：%d，错误内容：%s", apiResponse.StatusCode, apiResponse.StatusMessage)
		return client.TaskReportVO{}, flog.Error(log)
	}
	return apiResponse.Data, nil
}

func (receiver clientHttp) Kill(do client.DomainObject, taskId int64) error {
	clientUrl := fmt.Sprintf("http://%s:%d/api/kill", do.Ip, do.Port)
	var apiResponse core.ApiResponse[any]
	body := map[string]any{
		"TaskId": taskId,
	}
	_, err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(body).PostUnmarshal(&apiResponse)
	if err != nil {
		return err
	}
	if apiResponse.StatusCode != 200 {
		return flog.Errorf("状态码：%d，错误内容：%s", apiResponse.StatusCode, apiResponse.StatusMessage)
	}
	return err
}
