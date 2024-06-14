package response

type InfoResponse struct {
	TaskGroupCount      int64 // 任务组数量
	TaskGroupUnRunCount int   // 未运行任务组数量
	TodayFailCount      int64 // 今日失败任务数量
}
