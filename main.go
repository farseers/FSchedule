package main

import (
	"FSchedule/application/clientApp"
	"FSchedule/application/taskGroupApp"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/webapi"
	"net/http"
	"net/http/pprof"
)

func main() {
	fs.Initialize[StartupModule]("FSchedule")
	webapi.Area("/api/", func() {
		// 客户端注册
		webapi.RegisterPOST("/registry", clientApp.Registry)
		// 客户端下线
		webapi.RegisterPOST("/logout", clientApp.Logout)
		// 客户端回调
		webapi.RegisterPOST("/taskReport", taskGroupApp.TaskReport)
		// 上传日志
		webapi.RegisterPOST("/logReport", taskGroupApp.LogReport)
	})
	webapi.UseApiResponse()
	webapi.UsePprof()
	webapi.Run()
}

func initPprofMonitor() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	err := http.ListenAndServe(":8088", mux)
	if err != nil {
		_ = flog.Errorf("funcRetErr=http.ListenAndServe||err=%s", err.Error())
	}
}
