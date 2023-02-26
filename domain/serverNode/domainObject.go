package serverNode

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/parse"
	"strings"
	"time"
)

type DomainObject struct {
	Id         int64     // 客户端ID
	Name       string    // 客户端名称
	Ip         string    // 客户端IP
	Port       int       // 客户端端口
	ActivateAt time.Time // 活动时间
}

func New() *DomainObject {
	addr := configure.GetString("WebApi.Url")
	if addr == "" {
		addr = ":8888"
	}
	addr, _ = strings.CutPrefix(addr, ":")
	return &DomainObject{
		Id:         fs.AppId,
		Name:       fs.HostName,
		Ip:         fs.AppIp,
		Port:       parse.Convert(addr, 0),
		ActivateAt: time.Now(),
	}
}
