package repository

import (
	"FSchedule/domain/enum/executeStatus"
	"FSchedule/domain/enum/scheduleStatus"
	"FSchedule/domain/taskGroup"
	"FSchedule/infrastructure/repository/context"
	"FSchedule/infrastructure/repository/model"
	_ "embed"
	"sync"
	"time"

	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/redis"
)

var lock = &sync.Mutex{}

type taskRepository struct {
}

func getCacheManager(taskGroupName string) cache.ICacheManage[taskGroup.TaskEO] {
	key := "FSchedule_Task:" + taskGroupName
	if !container.IsRegister[cache.ICacheManage[taskGroup.TaskEO]](key) {
		lock.Lock()
		defer lock.Unlock()
		if !container.IsRegister[cache.ICacheManage[taskGroup.TaskEO]](key) {
			cacheManage := redis.SetProfiles[taskGroup.TaskEO](key, "Id", "default")
			cacheManage.SetItemSource(func(cacheId any) (taskGroup.TaskEO, bool) {
				po := context.MysqlContextIns("获取任务Task").Task.Where("Id = ?", cacheId).ToEntity()
				if po.Id > 0 {
					return mapper.Single[taskGroup.TaskEO](&po), true
				}
				return taskGroup.TaskEO{}, false
			})
		}
	}

	return container.Resolve[cache.ICacheManage[taskGroup.TaskEO]](key)
}

func (receiver *taskRepository) GetTask(taskGroupName string, taskId int64) taskGroup.TaskEO {
	item, _ := getCacheManager(taskGroupName).GetItem(taskId)
	return item
}

func (receiver *taskRepository) SaveTask(taskEO taskGroup.TaskEO) {
	getCacheManager(taskEO.Name).SaveItem(taskEO)
}

func (receiver *taskRepository) RemoveCache(taskGroupName string, lstSave collections.List[model.TaskPO]) {
	if lstSave.Count() == 0 {
		return
	}
	cacheManager := getCacheManager(taskGroupName)
	// 清除缓存
	lstSave.Foreach(func(po *model.TaskPO) {
		cacheManager.Remove(po.Id)
	})
}

// 获取要保存到数据库的任务列表
func (receiver *taskRepository) getSaveTaskList(taskGroupName string) collections.List[model.TaskPO] {
	cacheManager := getCacheManager(taskGroupName)
	lst := cacheManager.Get()
	lstSave := collections.NewList[model.TaskPO]()
	for i := 0; i < lst.Count(); i++ {
		do := lst.Index(i)
		// 保存成功后，已完成的任务，且最后运行时间大于1分钟的，移除列表
		// 最后运行时间超过1小时的移除。（如果有读取，还是会从数据库重新读的）
		if (do.IsFinish() && dateTime.Now().Sub(do.RunAt).Seconds() >= float64(30)) ||
			(dateTime.Now().Sub(do.RunAt).Hours() >= float64(1)) {
			po := mapper.Single[model.TaskPO](&do)
			if po.CreateAt.Year() < 2000 {
				po.CreateAt = time.Date(2000, 0, 1, 0, 0, 0, 0, time.Local)
			}
			if po.FinishAt.Year() < 2000 {
				po.FinishAt = time.Date(2000, 0, 1, 0, 0, 0, 0, time.Local)
			}
			if po.SchedulerAt.Year() < 2000 {
				po.SchedulerAt = time.Date(2000, 0, 1, 0, 0, 0, 0, time.Local)
			}
			lstSave.Add(po)
		}
	}

	return lstSave
}

func (receiver *taskRepository) DeleteTask(taskGroupName string) {
	_, _ = context.MysqlContextIns("删除任务Task").Task.Where("name = ?", taskGroupName).Delete()
	getCacheManager(taskGroupName).Clear()
}

func (receiver *taskRepository) ToHistoryTaskList(clientName, taskGroupName string, scheduleStatus scheduleStatus.Enum, executeStatus executeStatus.Enum, taskId string, pageSize int, pageIndex int) collections.PageList[taskGroup.TaskEO] {
	ts := context.MysqlContextIns("获取任务Task列表").Task.Desc("create_at")
	ts = ts.WhereIf(taskGroupName != "", "name = ?", taskGroupName)
	ts = ts.WhereIf(clientName != "", "client_name = ?", clientName)
	ts = ts.WhereIf(scheduleStatus > -1, "schedule_status = ?", scheduleStatus)
	ts = ts.WhereIf(executeStatus > -1, "execute_status = ?", executeStatus)
	ts = ts.WhereIf(taskId != "", "id = ?", parse.ToInt64(taskId))
	lstPO := ts.ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[taskGroup.TaskEO](lstPO)
}

func (receiver *taskRepository) ToTaskSpeedList() collections.List[taskGroup.TaskEO] {
	sql := "SELECT name, CAST(avg(`run_speed`) as UNSIGNED) as `run_speed` FROM `fschedule_task` WHERE execute_status = ? and create_at >= DATE_SUB(CURDATE(), INTERVAL 3 DAY) group by name"
	lstPO := context.MysqlContextIns("计算任务Task速度").Task.ExecuteSqlToList(sql, executeStatus.Success)
	return mapper.ToList[taskGroup.TaskEO](lstPO)
}

// TaskClearFinish 清除成功的任务记录（1天前）
func (receiver *taskRepository) TaskClearFinish(taskGroupName string, taskId int64) {
	_, _ = context.MysqlContextIns("清除成功的任务记录").Task.Where("name = ? and (execute_status = ? or execute_status = ?) and create_at < ? and Id < ?", taskGroupName, executeStatus.Success, executeStatus.Fail, time.Now().Add(-24*time.Hour), taskId).Delete()
}

func (receiver *taskRepository) ToTaskFinishList(taskGroupName string, top int) collections.List[taskGroup.TaskEO] {
	lstPO := context.MysqlContextIns("获取已完成的任务Task").Task.Where("name = ? and (execute_status = ? or execute_status = ?)", taskGroupName, executeStatus.Success, executeStatus.Fail).Desc("create_at").Limit(top).ToList()
	return mapper.ToList[taskGroup.TaskEO](lstPO)
}

// 获取已完成的任务TaskId
func (receiver *taskRepository) GetLastFinishTaskId(reservedTaskCount int) map[string]int64 {
	var m map[string]int64
	sql := `WITH ranked_tasks AS (
			SELECT
				id,
				name,
				ROW_NUMBER() OVER (PARTITION BY name ORDER BY create_at DESC) AS row_num
			FROM fschedule_task
			WHERE execute_status IN (2, 3)
			)
			SELECT name, id AS cutoff_id
			FROM ranked_tasks
			WHERE row_num = ?;`
	context.MysqlContextIns("获取已完成的任务TaskId").ExecuteSqlToMap(&m, sql, reservedTaskCount+1)
	return m
}

func (receiver *taskRepository) ToTaskFinishPageList(pageSize int, pageIndex int) collections.PageList[taskGroup.TaskEO] {
	page := context.MysqlContextIns("获取已完成的任务Task").Task.Where("(execute_status = ? or execute_status = ?) and (create_at >= ?)", executeStatus.Fail, executeStatus.Success, time.Now().Add(-24*time.Hour)).
		Desc("run_at").ToPageList(pageSize, pageIndex)
	return receiver.toTaskPageListTaskEO(page)
}

func (receiver *taskRepository) TodayFailCount() int64 {
	now := dateTime.Now()
	return context.MysqlContextIns("今日失败的任务数量").Task.Where("execute_status = ? and create_at >= ?", executeStatus.Fail, now.Date()).Count()
}

func (receiver *taskRepository) toTaskPageListTaskEO(page collections.PageList[model.TaskPO]) collections.PageList[taskGroup.TaskEO] {
	lst := mapper.ToList[taskGroup.TaskEO](page.List)
	return collections.NewPageList(lst, page.RecordCount)
}

func (receiver *taskRepository) GetStatCount() collections.List[taskGroup.StatTaskEO] {
	var array []taskGroup.StatTaskEO
	taskStatCountSql :=
		`SELECT client_name, execute_status, COUNT(*) AS count
		 FROM fschedule.fschedule_task
		 WHERE create_at >= (NOW() - INTERVAL 30 MINUTE) and client_name !=''
		 GROUP BY client_name,execute_status;`
	_, _ = context.MysqlContextIns("统计任务成功失败数量").ExecuteSqlToResult(&array, taskStatCountSql)
	return mapper.ToList[taskGroup.StatTaskEO](array)
}
