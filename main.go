package main

import (
	_ "bbs/routers"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
	"github.com/beego/beego/v2/task"
	"github.com/tidwall/gjson"
	_ "gorm.io/gorm"
	"os"
	"runtime"
	"strings"
	"time"
)

// Reply for AOP
type Reply struct {
	Time      int64       `json:"timestamp"`
	ReplyCode int         `json:"code"`
	Result    string      `json:"result"`
	Data      interface{} `json:"data"`
}

func init() {
	initLogs()
	printLogs()
	initConfig()
}

func main() {
	hostName, _ := os.Hostname()
	logs.Info("appname:", web.BConfig.AppName, "\n\tENV:", runtime.Version(), runtime.GOOS, runtime.GOARCH, hostName)
	logs.Debug("web.BConfig.RunMode =", web.BConfig.RunMode)
	if web.BConfig.RunMode == "dev" {
		web.BConfig.WebConfig.DirectoryIndex = true
		web.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	web.BConfig.WebConfig.StaticDir["/images"] = "static/images"
	web.BConfig.WebConfig.StaticDir["/live"] = "static/movies"
	web.BConfig.WebConfig.StaticDir["/"] = "static/web"
	// AOP
	regFilter()
	// task
	task.StartTask()
	// ListenAndServer()
	web.Run()
}

func regFilter() {
	// 验证传入的 POST 数据
	var DataParse = func(ctx *context.Context) {
		logs.Debug("[Filter] DataParse ", ctx.Request.Method, ctx.Request.RequestURI, ctx.Input.Method())
		if ctx.Input.Method() == "POST" {
			if !strings.Contains(ctx.Input.Header("Content-Type"), "application/json") {
				return
			}
			raw := ctx.Input.RequestBody
			if !gjson.Valid(string(raw)) {
				logs.Debug("[Filter] DataParse Valid json,check syntax :", string(raw))
				ctx.Output.Status = 415
				err := ctx.Output.JSON(&Reply{
					Time:      time.Now().Unix(),
					ReplyCode: 415,
					Result:    "Valid json,check syntax",
					Data:      nil,
				}, false, false)
				if err != nil {
					return
				}
			} else {
				timeNow := gjson.Get(string(raw), "timestamp")
				request := gjson.Get(string(raw), "request")
				logs.Debug("[Filter] DataParse timeNow:", timeNow, "request:", request)
				return
			}
		}
	}
	web.InsertFilter("*", web.BeforeExec, DataParse)
	web.InsertFilter("*", web.BeforeStatic, cors.Allow(&cors.Options{
		//AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		//ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		//AllowOrigins:     []string{"*","127.0.0.1:*"},// Access-Control-Allow-Origin
		AllowMethods: []string{"GET", "POST"},
		//AllowAllOrigins:  true,
		AllowCredentials: true,
	}))

}

func initLogs() {
	filename, _ := web.AppConfig.String("log.ProjLog")
	err := logs.SetLogger(logs.AdapterFile, `{"filename":"`+filename+`","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":7,"color":false}`)
	if err != nil || web.BConfig.RunMode == "dev" {
		_ = logs.SetLogger(logs.AdapterConsole)
		//logs.Error("unable to set logger to file: ",err.Error())
		logs.Info("set logs to console")
	}
	logs.Info("set logger to file: /project.log")
	logs.EnableFuncCallDepth(true)
	logs.Async(16)
}

func initConfig() {
	logs.Info("web.BConfig.CopyRequestBody: ", web.BConfig.CopyRequestBody)
	web.BConfig.CopyRequestBody = true
	RunMode := web.AppConfig.DefaultString("runmode", "dev")
	logs.Info("RunMode: ", RunMode)
}

// printLogs 测试日志
func printLogs() {
	logs.Emergency("logs.Emergency")
	logs.Alert("logs.Alert")       //1
	logs.Critical("logs.Critical") //2
	logs.Error("logs.Error")       //3
	logs.Warn("logs.Warning")      //4
	logs.Notice("logs.Notice")     //5
	logs.Info("logs.Info")         //6
	logs.Debug("logs.Debug")       //7
	logs.Trace("logs.Trace")
}
