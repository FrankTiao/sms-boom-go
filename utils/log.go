package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var logFilePath string

func InitLog() {
	log.SetPrefix("[version:" + Version + "]")

	daily := time.Now().Format("2006-01-02")
	if _, err := DirExistsOrCreate(GetAppDataLogDir(daily)); err != nil {
		log.Printf("日志目录创建失败，err：%v", err)
		return
	}

	logFilePath = GetAppDataLogDir(daily, time.Now().Format("15")+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Printf("日志文件创建失败，err：%v", err)
		return
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func PanicHandler(err any) {
	if err == nil {
		return
	}

	var buf [4096]byte
	n := runtime.Stack(buf[:], false)

	log.Printf("panic: %v, Stack: %v", err, string(buf[:n]))
	fmt.Printf("\n发生预期外的系统错误：%s \n"+
		"您可以按 Ctrl+C 键退出后重新打开应用\n"+
		"\n若仍然报错："+
		"请前往 https://github.com/franktiao/sms-boom-go/issues 反馈该问题 \n"+
		"反馈问题时请提供位于 %s 的日志文件 \n"+
		"开发者看到后会第一时间处理\n",
		err,
		strings.ReplaceAll(logFilePath, "\\", "/"),
	)
}
