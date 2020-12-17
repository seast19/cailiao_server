package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// 全局logger
var Mlogger *logrus.Logger

//设置logger
func SetLogger() {
	//日志名
	layoyt := "2006_01_02"
	fileName := fmt.Sprintf("./logs/log_%s.txt", time.Now().Format(layoyt))

	fmt.Println(fileName)

	//确保日志文件存在
	checkAndCreateLogFile(fileName)

	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("log error", err)
		return
	}

	//实例化
	Mlogger = logrus.New()

	//设置输出
	Mlogger.Out = src

	//设置日志级别
	Mlogger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	//Mlogger.SetFormatter(&logrus.TextFormatter{})
	//设置日志格式
	Mlogger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	//从通道读取数据写入日志文件
	Mlogger.Info("+----- 日志系统加载成功 -----+")

}

//检查日志文件是否存在，不在存在则创建
func checkAndCreateLogFile(fn string) {
	//	无该文件则创建
	_, err := os.Stat(fn)
	if err == nil {
		return
	}

	fmt.Println("日志文件不存在，创建日志文件")

	fs, err := os.Create(fn)
	if err != nil {
		panic(err)
	}

	defer fs.Close()

}
