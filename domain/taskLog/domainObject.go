package taskLog

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/dateTime"
	"time"
)

type DomainObject struct {
	Id          int64                                  // 主键
	TaskGroupId int64                                  // 任务组ID
	Name        string                                 // 实现Job的特性名称（客户端识别哪个实现类）
	Caption     string                                 // 任务组标题
	Ver         int                                    // 版本
	TaskId      int64                                  // 任务ID
	Data        collections.Dictionary[string, string] // 本次执行任务时的Data数据
	LogLevel    eumLogLevel.Enum                       // 日志级别
	Content     string                                 // 日志内容
	CreateAt    dateTime.DateTime                      // 日志时间
}

func NewDO(mame, caption string, ver int, taskId int64, taskGroupId int64, data collections.Dictionary[string, string], logLevel eumLogLevel.Enum, content string, createAt int64) DomainObject {
	do := DomainObject{
		Id:          0,
		TaskGroupId: taskGroupId,
		Name:        mame,
		Caption:     caption,
		Ver:         ver,
		TaskId:      taskId,
		Data:        data,
		LogLevel:    logLevel,
		Content:     content,
		CreateAt:    dateTime.New(time.UnixMilli(createAt)),
	}
	return do
}
