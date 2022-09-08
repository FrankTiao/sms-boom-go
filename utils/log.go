package utils

import (
	"log"
	"os"
	"time"
)

func InitLog() {
	log.SetPrefix("[version:" + Version + "]")

	daily := time.Now().Format("2006-01-02")
	if _, err := DirExistsOrCreate(GetAppDataLogDir(daily)); err != nil {
		log.Printf("日志目录创建失败，err：%v", err)
		return
	}

	file := GetAppDataLogDir(daily, time.Now().Format("15")+".log")
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Printf("日志文件创建失败，err：%v", err)
		return
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
