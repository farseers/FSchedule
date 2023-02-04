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

func (receiver clientHttp) Check(do *client.DomainObject) bool {
	clientUrl := fmt.Sprintf("http://%s:%d/api/check", do.Ip, do.Port)
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).PostUnmarshal(&apiResponse)
	if err != nil {
		return false
	}
	if apiResponse.StatusCode != 200 {
		flog.Infof("客户端：http://%s:%d，状态码：%s", do.Ip, do.Port, apiResponse.StatusCode)
		return false
	}
	return true
}

func (receiver clientHttp) Invoke(do *client.DomainObject, task *client.TaskEO) bool {
	clientUrl := fmt.Sprintf("http://%s:%d/api/invoke", do.Ip, do.Port)
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).Body(task).PostUnmarshal(&apiResponse)
	if err != nil {
		return false
	}
	if apiResponse.StatusCode != 200 {
		flog.Infof("客户端：http://%s:%d，状态码：%s", do.Ip, do.Port, apiResponse.StatusCode)
		return false
	}
	return true
}

func (receiver clientHttp) Status(do *client.DomainObject) bool {
	clientUrl := fmt.Sprintf("http://%s:%d/api/status", do.Ip, do.Port)
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).PostUnmarshal(&apiResponse)
	if err != nil {
		return false
	}
	if apiResponse.StatusCode != 200 {
		flog.Infof("客户端：http://%s:%d，状态码：%s", do.Ip, do.Port, apiResponse.StatusCode)
		return false
	}
	return true
}

func (receiver clientHttp) Kill(do *client.DomainObject) bool {
	clientUrl := fmt.Sprintf("http://%s:%d/api/kill", do.Ip, do.Port)
	var apiResponse core.ApiResponse[any]
	err := http.NewClient(clientUrl).HeadAdd(tokenName, token).PostUnmarshal(&apiResponse)
	if err != nil {
		return false
	}
	if apiResponse.StatusCode != 200 {
		flog.Infof("客户端：http://%s:%d，状态码：%s", do.Ip, do.Port, apiResponse.StatusCode)
		return false
	}
	return true
}
