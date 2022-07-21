package utils

import (
	"github.com/beego/beego/v2/core/logs"
)

// SetLogger2 beego的log配置
func SetLogger() {
	//an official log.Logger
	//l := logs.GetLogger()

	//输出console
	//_ = logs.SetLogger(logs.AdapterConsole)

	//输出file
	_ = logs.SetLogger(logs.AdapterFile, `{"filename":"logs/logs.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":15,"color":true}`)

	//l.Println("this is a message of http")
	//an official log.Logger with prefix ORM
	//logs.GetLogger("ORM").Println("this is a message of orm")

	//logs.Debug("my book is bought in the year of ", 2016)
	//logs.Info("this %s cat is %v years old", "yellow", 3)
	//logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
	//logs.Error(1024, "is a very", "good game")
	//logs.Critical("oh,crash")
}
